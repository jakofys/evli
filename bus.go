package evli

import (
	"reflect"
	"runtime"
	"sync"
	"time"
)

// Bus is pool that bindinb event to his listeners
type Bus struct {
	routing map[string][]listenerDesc
}

// broadcaster is binding between event and multiple listener to this event
var broadcaster = NewBus()

// NewBus
func NewBus() *Bus {
	return &Bus{routing: make(map[string][]listenerDesc)}
}

// Register a listener and binding it to a specific event
func (b *Bus) Listen(e Event, listeners ...Listener) error {
	name := reflect.TypeOf(e).Name()
	if _, ok := b.routing[name]; !ok {
		b.routing[name] = []listenerDesc{}
	}
	for _, listener := range listeners {
		ln := runtime.FuncForPC(reflect.ValueOf(listener).Pointer()).Name()
		if listener == nil {
			return ErrInvalidListener
		}
		for _, l := range b.routing[name] {
			if ln == l.name {
				return ErrDuplicatedListener
			}
		}
		b.routing[name] = append(b.routing[name], listenerDesc{name: ln, l: listener})
	}
	return nil
}

// Emit an event through the all listeners
// return ErrNoListenerFound if no listener register to this event
func (b *Bus) Emit(e Event, m Meta) error {
	name := reflect.TypeOf(e).Name()
	k := reflect.ValueOf(e).Kind()
	if k == reflect.Ptr {
		return ErrNotPointerEvent
	}
	if listeners, ok := b.routing[name]; ok {
		if len(listeners) == 0 {
			return ErrNoListenerFound
		}
		s := &Source{
			payload:   e,
			emittedAd: time.Now(),
			schema:    schema(e),
			meta:      m,
			mu:        sync.Mutex{},
		}

		for _, listener := range listeners {
			go listener.l(s)
		}
	} else {
		return ErrUnknownEvent
	}
	return nil
}

// Listen a listener and binding it to a specific event
func Listen(e Event, listeners ...Listener) error {
	return broadcaster.Listen(e, listeners...)
}

// Emit an event through the all listeners
// return ErrNoListenerFound if no listener register to this event
func Emit(e Event, m Meta) error {
	return broadcaster.Emit(e, m)
}
