package dispatcher

import (
	"github.com/gopherd/doge/event"
)

var _ event.Dispatcher[any] = (*Dispatcher[any])(nil)

// Dispatcher implements event.Dispatcher
type Dispatcher[T comparable] struct {
	listeners map[T][]event.Listener[T]
}

// AddEventListener implements event.Dispatcher AddEventListener method
func (dispatcher *Dispatcher[T]) AddEventListener(listener event.Listener[T]) {
	if dispatcher.listeners == nil {
		dispatcher.listeners = make(map[T][]event.Listener[T])
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

// HasEventListener implements event.Dispatcher HasEventListener method
func (dispatcher *Dispatcher[T]) HasEventListener(listener event.Listener[T]) bool {
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

// RemoveEventListener implements event.Dispatcher RemoveEventListener method
func (dispatcher *Dispatcher[T]) RemoveEventListener(listener event.Listener[T]) bool {
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

// DispatchEvent implements event.Dispatcher DispatchEvent method
func (dispatcher *Dispatcher[T]) DispatchEvent(event event.Event[T]) bool {
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
