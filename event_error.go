package evli

import "errors"

var (
	// ErrNoListenerFound if no listener found when emitting event
	ErrNoListenerFound = errors.New("evli: no listener found when emitting event")

	// ErrInvalidListener prevent for avoid nil func
	ErrInvalidListener = errors.New("evli: listener must implement Listener interface")

	// ErrUnknownEvent any event declared with a listener
	ErrUnknownEvent = errors.New("evli: event haven't be declared in register function")

	// ErrDuplicatedListener when trying to declare multiple listeners which are the same
	ErrDuplicatedListener = errors.New("evli: duplicated listeners are register for an event")

	// ErrNotPointerEvent only none pointer can be used
	ErrNotPointerEvent = errors.New("evli: event to emitt can't be a pointer")
)
