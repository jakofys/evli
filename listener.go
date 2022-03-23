package evli

// Subscriber listen to an event and be execute
// context if used to canceled event emittion process in gross case.
type Subscriber func(*Source)

// subscriberDesc is the descriptor of subscriber
type subscriberDesc struct {
	l    Subscriber
	name string
}
