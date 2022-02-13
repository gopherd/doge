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

// Dispatcher represents an event dispatcher
type Dispatcher[T comparable] interface {
	// AddEventListener registers a Listener and returns the listener ID
	AddEventListener(listener Listener[T]) (id int)
	// HasEventListener reports whether the Dispatcher has specified listener
	HasEventListener(id int) bool
	// RemoveEventListener removes specified listener
	RemoveEventListener(id int) bool
	// DispatchEvent dispatchs event
	DispatchEvent(event Event[T]) bool
}
