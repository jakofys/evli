package evli_test

import (
	"testing"
	"time"

	"github.com/jakofys/evli"
)

type EventStd struct {
	Name   string
	Age    int
	hidden string
}

type EventRandom struct {
	Name   string
	Number int
}

type EventNoListener struct {
	Name string
}
type UnknownEvent struct {
	Name string
}

type RandomTaggedData struct {
	UserName string `field:"Name"`
	Age      int
}

func EventRandomListener(s *evli.Source) {
	m := s.Meta()
	t := m["testing"].(*testing.T)
	e := EventRandom{}
	if s.EventName(evli.ShortName) != "EventStd" {
		t.Errorf("expected EventStd and have %s", s.EventName(evli.ShortName))
	}
	if s.EventName(evli.LongName) != "evli_test.EventStd" {
		t.Errorf("expected evli_test.EventStd and have %s", s.EventName(evli.LongName))
	}
	if s.EventName(evli.PackageName) != "github.com/jakofys/evli/evli_test.EventStd" {
		t.Errorf("expected  github.com/jakofys/evli/evli_test.EventStd and have %s", s.EventName(evli.PackageName))
	}
	if time.Now().Before(s.EmittedAt()) {
		t.Errorf("expected date emitted before now registered")
	}

	if err := s.Payload(&e); err != nil {
		t.Errorf("expected Listen to be operationnal: %s", err)
	}
	if e.Name != "My name" && e.Number != 0 {
		t.Errorf("expected not listening public data correctly")
	}
}

func EventStdListener(s *evli.Source) {
	m := s.Meta()
	t := m["testing"].(*testing.T)
	e := EventStd{}
	if s.EventName(evli.ShortName) != "EventStd" {
		t.Errorf("expected EventStd and have %s", s.EventName(evli.ShortName))
	}
	if err := s.Payload(&e); err != nil {
		t.Errorf("expected Payload to be operationnal: %s", err)
	}
	if e.Name != "My name" && e.Age != 1 {
		t.Errorf("expected not listening public data correctly")
	}

	obj := &RandomTaggedData{}
	if err := s.Payload(obj); err != nil {
		t.Errorf("expected Payload to be operationnal: %s", err)
	}
	if e.Age != 1 {
		t.Errorf("expected not listening public data correctly")
	}
	if obj.UserName != "My name" {
		t.Errorf("expected not listening public data correctly")
	}
}

func TestRegisterListeners(t *testing.T) {
	type Give struct {
		listener []evli.Listener
		event    evli.Event
	}
	table := []struct {
		give    Give
		expect  error
		message string
	}{
		{
			give: Give{
				listener: []evli.Listener{EventStdListener, EventRandomListener},
				event:    EventStd{},
			},
			expect:  nil,
			message: "don't register basic listener",
		},
		{
			give: Give{
				listener: []evli.Listener{EventStdListener, EventStdListener},
				event:    EventStd{},
			},
			expect:  evli.ErrDuplicatedListener,
			message: "don't detect duplicated listener",
		},
		{
			give: Give{
				listener: []evli.Listener{nil},
				event:    EventStd{},
			},
			expect:  evli.ErrInvalidListener,
			message: "don't detect invalid listener as nil",
		},
		{
			give: Give{
				listener: []evli.Listener{},
				event:    EventNoListener{},
			},
			expect:  nil,
			message: "must save event name",
		},
	}
	for _, test := range table {
		if err := evli.Listen(test.give.event, test.give.listener...); err != test.expect {
			t.Errorf(test.message, err)
		}
	}
}

func TestEmittingEvent(t *testing.T) {
	table := []struct {
		give    evli.Event
		expect  error
		message string
	}{
		{
			give:    EventStd{Name: "My name", Age: 1, hidden: "hi hi hi"},
			expect:  nil,
			message: "don't emit basic listener event",
		},
		{
			give:    &EventStd{Name: "My name", Age: 1, hidden: "hi hi hi"},
			expect:  evli.ErrNotPointerEvent,
			message: "does not support pointer event",
		},
		{
			give:    EventNoListener{},
			expect:  evli.ErrNoListenerFound,
			message: "don't treat zero listener",
		},
		{
			give:    UnknownEvent{},
			expect:  evli.ErrUnknownEvent,
			message: "don't detect unknown event",
		},
	}
	m := map[string]interface{}{"testing": t}
	for _, test := range table {
		if err := evli.Emit(test.give, m); err != test.expect {
			t.Errorf(test.message, err)
		}
	}
}
