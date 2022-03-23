package evli

import (
	"reflect"
	"runtime"
	"sync"
	"time"
)

// Bus is pool that bindinb event to his Subscribers
type Bus struct {
	routing map[string][]subscriberDesc
}

// broadcaster is binding between event and multiple Subscriber to this event
var broadcaster = NewBus()

// NewBus
func NewBus() *Bus {
	return &Bus{routing: make(map[string][]subscriberDesc)}
}

// Register a Subscriber and binding it to a specific event
func (b *Bus) Subscribe(e Event, Subscribers ...Subscriber) error {
	name := reflect.TypeOf(e).Name()
	if _, ok := b.routing[name]; !ok {
		b.routing[name] = []subscriberDesc{}
	}
	for _, Subscriber := range Subscribers {
		ln := runtime.FuncForPC(reflect.ValueOf(Subscriber).Pointer()).Name()
		if Subscriber == nil {
			return ErrInvalidSubscriber
		}
		for _, l := range b.routing[name] {
			if ln == l.name {
				return ErrDuplicatedSubscriber
			}
		}
		b.routing[name] = append(b.routing[name], subscriberDesc{name: ln, l: Subscriber})
	}
	return nil
}

// Emit an event through the all Subscribers
// return ErrNoSubscriberFound if no Subscriber register to this event
func (b *Bus) Emit(e Event, m Meta) error {
	name := reflect.TypeOf(e).Name()
	k := reflect.ValueOf(e).Kind()
	if k == reflect.Ptr {
		return ErrNotPointerEvent
	}
	if Subscribers, ok := b.routing[name]; ok {
		if len(Subscribers) == 0 {
			return ErrNoSubscriberFound
		}
		s := &Source{
			payload:   e,
			emittedAd: time.Now(),
			schema:    schema(e),
			meta:      m,
			mu:        sync.Mutex{},
		}

		for _, Subscriber := range Subscribers {
			go Subscriber.l(s)
		}
	} else {
		return ErrUnknownEvent
	}
	return nil
}

// Subscribe a Subscriber and binding it to a specific event
func Subscribe(e Event, Subscribers ...Subscriber) error {
	return broadcaster.Subscribe(e, Subscribers...)
}

// Emit an event through the all Subscribers
// return ErrNoSubscriberFound if no Subscriber register to this event
func Emit(e Event, m Meta) error {
	return broadcaster.Emit(e, m)
}
