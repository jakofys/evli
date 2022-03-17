package evli

// Listener listen to an event and be execute
// context if used to canceled event emittion process in gross case.
type Listener func(*Source)

// listenerDesc is the descriptor of listener
type listenerDesc struct {
	l    Listener
	name string
}
