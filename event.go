package evli

// Event is a structure that describe an declared event
// it must named has a passed event as "UserUpdated" and not "UserUpdate"
// recommanded to be a struct for clearly fields definition
// in func case, make sure that can support as a func type and not a random function, can panic
type Event = interface{}
