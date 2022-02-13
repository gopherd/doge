package dispatcher

import (
	"github.com/gopherd/doge/container"
	"github.com/gopherd/doge/event"
)

var _ event.Dispatcher[any] = (*Dispatcher[any])(nil)

// Dispatcher implements event.Dispatcher
type Dispatcher[T comparable] struct {
	nextId    event.ListenerID
	listeners map[T][]container.Pair[event.ListenerID, event.Listener[T]]
	mapping   map[event.ListenerID]container.Pair[T, int]
}

// AddEventListener implements event.Dispatcher AddEventListener method
func (dispatcher *Dispatcher[T]) AddEventListener(listener event.Listener[T]) event.ListenerID {
	if dispatcher.listeners == nil {
		dispatcher.listeners = make(map[T][]container.Pair[event.ListenerID, event.Listener[T]])
		dispatcher.mapping = make(map[event.ListenerID]container.Pair[T, int])
	}
	dispatcher.nextId++
	var id = dispatcher.nextId
	var eventType = listener.EventType()
	var listeners = dispatcher.listeners[eventType]
	var index = len(listeners)
	dispatcher.listeners[eventType] = append(listeners, container.MakePair(id, listener))
	dispatcher.mapping[id] = container.MakePair(eventType, index)
	return id
}

// HasEventListener implements event.Dispatcher HasEventListener method
func (dispatcher *Dispatcher[T]) HasEventListener(id event.ListenerID) bool {
	if dispatcher.mapping == nil {
		return false
	}
	_, ok := dispatcher.mapping[id]
	return ok
}

// RemoveEventListener implements event.Dispatcher RemoveEventListener method
func (dispatcher *Dispatcher[T]) RemoveEventListener(id event.ListenerID) bool {
	if dispatcher.mapping == nil {
		return false
	}
	var index, ok = dispatcher.mapping[id]
	if !ok {
		return false
	}
	delete(dispatcher.mapping, id)
	var listeners = dispatcher.listeners[index.First]
	var n = len(listeners) - 1
	if index.Second != n {
		listeners[index.Second] = listeners[n]
		dispatcher.mapping[listeners[index.Second].First] = index
	}
	listeners[n].Second = nil
	dispatcher.listeners[index.First] = listeners[:n]
	return true
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
		listeners[i].Second.Handle(event)
	}
	return true
}
