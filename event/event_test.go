package event_test

import (
	"testing"

	"github.com/gopherd/doge/event"
)

type testEvent struct {
}

func (e testEvent) Type() string {
	return "test"
}

func TestDispatchEvent(t *testing.T) {
	var fired bool
	dispatcher := event.NewDispatcher()
	dispatcher.AddEventListener("test", event.Listen(func(e testEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(testEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

func TestDispatchEventPointer(t *testing.T) {
	var fired bool
	dispatcher := event.NewDispatcher()
	dispatcher.AddEventListener("test", event.Listen(func(e *testEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(new(testEvent))
	if !fired {
		t.Fatal("event not fired")
	}
}
