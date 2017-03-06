package msgrouter

// Component is an interface for go routines which will be routed from or to
//
// Send is the interface in which the router will use to send a message to
// this go routine. This is usually a wrapper around the go routines IN channel.
//
// SetID and GetID will be used to register and lookup our components in the
// router
type Component interface {
	Send(interface{}) error
	SetID(ComponentID) error
	// TODO Determine best way to handle empty UUID.
	GetID() (ComponentID, error)
}
