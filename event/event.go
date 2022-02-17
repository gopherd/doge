package event

import (
	"fmt"

	"github.com/gopherd/doge/container"
)

// Event is the interface that wraps the basic Type method.
type Event[T comparable] interface {
	Type() T // Type gets type of event
}

// A Listener handles fired event
type Listener[T comparable] interface {
	EventType() T    // EventType gets type of listening event
	Handle(Event[T]) // Handle handles fired event
}

// Listen creates a Listener by eventType and handler function
func Listen[T comparable, E Event[T]](eventType T, handler func(E)) Listener[T] {
	return listenerFunc[T, E]{eventType, handler}
}

type listenerFunc[T comparable, E Event[T]] struct {
	eventType T
	handler   func(E)
}

func (h listenerFunc[T, E]) EventType() T {
	return h.eventType
}

func (h listenerFunc[T, E]) Handle(event Event[T]) {
	if e, ok := event.(E); ok {
		h.handler(e)
	} else {
		panic(fmt.Sprintf("unexpected event %T for type %v", event, event.Type()))
	}
}

// Dispatcher manages event listeners
type Dispatcher[T comparable] struct {
	nextid    int
	ordered   bool
	listeners map[T][]container.Pair[int, Listener[T]]
	mapping   map[int]container.Pair[T, int]
}

// Ordered reports whether the listeners fired by added order
func (dispatcher *Dispatcher[T]) Ordered() bool {
	return dispatcher.ordered
}

// SetOrdered sets whether the listeners fired by added order
func (dispatcher *Dispatcher[T]) SetOrdered(ordered bool) {
	dispatcher.ordered = ordered
}

// AddEventListener registers a Listener
func (dispatcher *Dispatcher[T]) AddEventListener(listener Listener[T]) int {
	if dispatcher.listeners == nil {
		dispatcher.listeners = make(map[T][]container.Pair[int, Listener[T]])
		dispatcher.mapping = make(map[int]container.Pair[T, int])
	}
	dispatcher.nextid++
	var id = dispatcher.nextid
	var eventType = listener.EventType()
	var listeners = dispatcher.listeners[eventType]
	var index = len(listeners)
	dispatcher.listeners[eventType] = append(listeners, container.MakePair(id, listener))
	dispatcher.mapping[id] = container.MakePair(eventType, index)
	return id
}

// HasEventListener reports whether the Dispatcher has specified listener
func (dispatcher *Dispatcher[T]) HasEventListener(id int) bool {
	if dispatcher.mapping == nil {
		return false
	}
	_, ok := dispatcher.mapping[id]
	return ok
}

// RemoveEventListener removes specified listener
func (dispatcher *Dispatcher[T]) RemoveEventListener(id int) bool {
	if dispatcher.listeners == nil {
		return false
	}
	index, ok := dispatcher.mapping[id]
	if !ok {
		return false
	}
	var eventType = index.First
	var listeners = dispatcher.listeners[eventType]
	var last = len(listeners) - 1
	if index.Second != last {
		if dispatcher.ordered {
			copy(listeners[index.Second:last], listeners[index.Second+1:])
			for i := index.Second; i < last; i++ {
				dispatcher.mapping[listeners[i].First] = container.MakePair(eventType, i)
			}
		} else {
			listeners[index.Second] = listeners[last]
			dispatcher.mapping[listeners[index.Second].First] = container.MakePair(eventType, index.Second)
		}
	}
	listeners[last].Second = nil
	dispatcher.listeners[eventType] = listeners[:last]
	delete(dispatcher.mapping, id)
	return true
}

// DispatchEvent dispatchs event
func (dispatcher *Dispatcher[T]) DispatchEvent(event Event[T]) bool {
	if dispatcher.listeners == nil {
		return false
	}
	listeners, ok := dispatcher.listeners[event.Type()]
	if !ok || len(listeners) == 0 {
		return false
	}
	for i := range listeners {
		listeners[i].Second.Handle(event)
	}
	return true
}
