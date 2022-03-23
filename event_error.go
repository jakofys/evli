package evli

import "errors"

var (
	// ErrNoSubscriberFound if no Subscriber found when emitting event
	ErrNoSubscriberFound = errors.New("evli: no Subscriber found when emitting event")

	// ErrInvalidSubscriber prevent for avoid nil func
	ErrInvalidSubscriber = errors.New("evli: Subscriber must implement Subscriber interface")

	// ErrUnknownEvent any event declared with a Subscriber
	ErrUnknownEvent = errors.New("evli: event haven't be declared in register function")

	// ErrDuplicatedSubscriber when trying to declare multiple Subscribers which are the same
	ErrDuplicatedSubscriber = errors.New("evli: duplicated Subscribers are register for an event")

	// ErrNotPointerEvent only none pointer can be used
	ErrNotPointerEvent = errors.New("evli: event to emitt can't be a pointer")
)
