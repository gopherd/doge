package event

import (
	"fmt"
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
	listeners map[T][]Listener[T]
}

// AddEventListener registers a Listener
func (dispatcher *Dispatcher[T]) AddEventListener(listener Listener[T]) {
	if dispatcher.listeners == nil {
		dispatcher.listeners = make(map[T][]Listener[T])
	}
	var eventType = listener.EventType()
	var listeners = dispatcher.listeners[eventType]
	for i := range listeners {
		if listeners[i] == listener {
			return
		}
	}
	dispatcher.listeners[eventType] = append(listeners, listener)
}

// HasEventListener reports whether the Dispatcher has specified listener
func (dispatcher *Dispatcher[T]) HasEventListener(listener Listener[T]) bool {
	if dispatcher.listeners == nil {
		return false
	}
	var listeners, ok = dispatcher.listeners[listener.EventType()]
	if !ok {
		return false
	}
	for i := range listeners {
		if listeners[i] == listener {
			return true
		}
	}
	return false
}

// RemoveEventListener removes specified listener
func (dispatcher *Dispatcher[T]) RemoveEventListener(listener Listener[T]) bool {
	if dispatcher.listeners == nil {
		return false
	}
	var eventType = listener.EventType()
	var listeners, ok = dispatcher.listeners[eventType]
	if !ok {
		return false
	}
	for i := range listeners {
		if listeners[i] == listener {
			var n = len(listeners)
			copy(listeners[i:n-1], listeners[i+1:])
			listeners[n-1] = nil
			dispatcher.listeners[eventType] = listeners
			return true
		}
	}
	return false
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
		listeners[i].Handle(event)
	}
	return true
}
