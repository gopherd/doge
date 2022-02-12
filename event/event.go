package event

import (
	"fmt"

	"github.com/gopherd/doge/container"
)

// ID represents type of event ID
type ID = int64

// Type represents type of event type
type Type = string

// Event is the interface that wraps the basic Type method.
type Event interface {
	Type() Type
}

// A Listener handles fired event
type Listener interface {
	Listen(Event)
}

// Listen wraps the handler func as a Listener
func Listen[E Event](handler func(E)) Listener {
	return listenerFunc[E]{handler}
}

type listenerFunc[E Event] struct {
	fn func(E)
}

func (h listenerFunc[E]) Listen(event Event) {
	if e, ok := event.(E); ok {
		h.fn(e)
	} else {
		panic(fmt.Sprintf("unexpected event %T for type %s", event, event.Type()))
	}
}

// Dispatcher represents an event dispatcher
type Dispatcher interface {
	// AddEventListener registers a Listener by Type and returns the event ID
	AddEventListener(Type, Listener) ID
	// HasEventListener reports whether the Dispatcher has specified event handler
	HasEventListener(Type, ID) bool
	// RemoveEventListener removes specified event handler
	RemoveEventListener(Type, ID) bool
	// DispatchEvent dispatchs event
	DispatchEvent(Event)
}

// BasicDispatcher implements a basic Dispatcher
type BasicDispatcher struct {
	nextId   ID
	handlers map[Type][]container.Pair[ID, Listener]
}

// NewDispatcher creates a BasicDispatcher
func NewDispatcher() *BasicDispatcher {
	return &BasicDispatcher{
		handlers: make(map[Type][]container.Pair[ID, Listener]),
	}
}

// AddEventListener implements Dispatcher AddEventListener method
func (dispatcher *BasicDispatcher) AddEventListener(eventType Type, handler Listener) ID {
	dispatcher.nextId++
	id := dispatcher.nextId
	dispatcher.handlers[eventType] = append(dispatcher.handlers[eventType], container.MakePair(id, handler))
	return id
}

// HasEventListener implements Dispatcher HasEventListener method
func (dispatcher *BasicDispatcher) HasEventListener(eventType Type, id ID) bool {
	if handlers, ok := dispatcher.handlers[eventType]; ok {
		for i := range handlers {
			if handlers[i].First == id {
				return true
			}
		}
	}
	return false
}

// RemoveEventListener implements Dispatcher RemoveEventListener method
func (dispatcher *BasicDispatcher) RemoveEventListener(eventType Type, id ID) bool {
	if handlers, ok := dispatcher.handlers[eventType]; ok {
		for i := range handlers {
			if handlers[i].First == id {
				copy(handlers[i:], handlers[i+1:])
				var n = len(handlers) - 1
				handlers[n].Second = nil
				dispatcher.handlers[eventType] = handlers[:n]
				return true
			}
		}
	}
	return false
}

// DispatchEvent implements Dispatcher DispatchEvent method
func (dispatcher *BasicDispatcher) DispatchEvent(event Event) bool {
	handlers, ok := dispatcher.handlers[event.Type()]
	if !ok || len(handlers) == 0 {
		return false
	}
	for i := range handlers {
		handlers[i].Second.Listen(event)
	}
	return true
}
