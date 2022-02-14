package event_test

import (
	"reflect"
	"testing"

	"github.com/gopherd/doge/event"
)

type testStringEvent struct {
}

func (e testStringEvent) Type() string {
	return "test"
}

func TestDispatchEvent(t *testing.T) {
	var fired bool
	var dispatcher event.Dispatcher[string]
	dispatcher.AddEventListener(event.Listen("test", func(e testStringEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(testStringEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

func TestDispatchEventPointer(t *testing.T) {
	var fired bool
	var dispatcher event.Dispatcher[string]
	dispatcher.AddEventListener(event.Listen("test", func(e *testStringEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(&testStringEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}

type testEvent struct {
}

var testEventType = reflect.TypeOf((*testEvent)(nil)).Elem()

func (e testEvent) Type() reflect.Type {
	return testEventType
}

func TestDispatchReflectEvent(t *testing.T) {
	var fired bool
	var dispatcher event.Dispatcher[reflect.Type]
	dispatcher.AddEventListener(event.Listen(testEventType, func(e *testEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(new(testEvent))
	if !fired {
		t.Fatal("event not fired")
	}
}

type testIntEvent struct {
}

func (e testIntEvent) Type() int {
	return 1
}

func TestDispatchIntEvent(t *testing.T) {
	var fired bool
	var dispatcher event.Dispatcher[int]
	dispatcher.AddEventListener(event.Listen(1, func(e testIntEvent) {
		fired = true
	}))
	dispatcher.DispatchEvent(testIntEvent{})
	if !fired {
		t.Fatal("event not fired")
	}
}
