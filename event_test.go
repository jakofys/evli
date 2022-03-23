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

type EventNoSubscriber struct {
	Name string
}
type UnknownEvent struct {
	Name string
}

type RandomTaggedData struct {
	UserName string `field:"Name"`
	Age      int
}

func EventRandomSubscriber(s *evli.Source) {
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
		t.Errorf("expected Subscribe to be operationnal: %s", err)
	}
	if e.Name != "My name" && e.Number != 0 {
		t.Errorf("expected not listening public data correctly")
	}
}

func EventStdSubscriber(s *evli.Source) {
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

func TestRegisterSubscribers(t *testing.T) {
	type Give struct {
		Subscriber []evli.Subscriber
		event      evli.Event
	}
	table := []struct {
		give    Give
		expect  error
		message string
	}{
		{
			give: Give{
				Subscriber: []evli.Subscriber{EventStdSubscriber, EventRandomSubscriber},
				event:      EventStd{},
			},
			expect:  nil,
			message: "don't register basic Subscriber",
		},
		{
			give: Give{
				Subscriber: []evli.Subscriber{EventStdSubscriber, EventStdSubscriber},
				event:      EventStd{},
			},
			expect:  evli.ErrDuplicatedSubscriber,
			message: "don't detect duplicated Subscriber",
		},
		{
			give: Give{
				Subscriber: []evli.Subscriber{nil},
				event:      EventStd{},
			},
			expect:  evli.ErrInvalidSubscriber,
			message: "don't detect invalid Subscriber as nil",
		},
		{
			give: Give{
				Subscriber: []evli.Subscriber{},
				event:      EventNoSubscriber{},
			},
			expect:  nil,
			message: "must save event name",
		},
	}
	for _, test := range table {
		if err := evli.Subscribe(test.give.event, test.give.Subscriber...); err != test.expect {
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
			message: "don't emit basic Subscriber event",
		},
		{
			give:    &EventStd{Name: "My name", Age: 1, hidden: "hi hi hi"},
			expect:  evli.ErrNotPointerEvent,
			message: "does not support pointer event",
		},
		{
			give:    EventNoSubscriber{},
			expect:  evli.ErrNoSubscriberFound,
			message: "don't treat zero Subscriber",
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
