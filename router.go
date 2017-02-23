package msgrouter

import "errors"

// Router allows synchronization and routing decisions to be made between
// go routines. By having go routines send to a router we can create different
// messaging topologies. Operations on the router (route add, route remove,
// registering components, etc...) should block the receiving of messages from
// external go routines thus synchronizing updates without the need for locks
type Router interface {
	Send(msg *interface{}) error
	RegisterComponent(c Component) error
	UnregisterComponent(c Component) error
	AddRoute(src Component, dst Component) error
	RemoveRoute(src Component, dst Component) error
	ListRoutes() (string, error)
	Consume()
}

// componentID is an ID used to select registered components
type componentID int

// Map which correlates source component to one or more destination components
type routingTable map[Component][]Component

// GenericRouter is an implementation of a router. External channels are for
// API access while internal channels are for consuming off of.
// TODO: make struct or interface to handle messages on each channel, removing
// interface{}
type GenericRouter struct {
	externalMsgChan chan<- interface{}
	internalMsgChan <-chan interface{}
	externalRtChan  chan<- interface{}
	internalRtChan  <-chan interface{}
	externalRegChan chan<- interface{}
	internalRegChan <-chan interface{}
	rt              routingTable
	rc              []Component
}

// msg* structs are used to package messages that will be sent on the
// associated channel. External API
type msgMsg struct {
	src     componentID
	payload interface{}
}

type msgRt struct {
	src  componentID
	dest componentID
}

type msgReg struct {
	c Component
}

// NewGenericRouter is a constructor for a generic implementation of a Router
// Channels should be buffered so that sending go routines do not block while
// blocking operations occur on router
func NewGenericRouter(bufferSize int) *GenericRouter {

	// make channels
	msgChan := make(chan interface{}, bufferSize)
	rtChan := make(chan interface{}, bufferSize)
	cmpChan := make(chan interface{}, bufferSize)

	// create routing table
	rt := routingTable{}

	// create registeredComponents array
	rc := []Component{}

	// construct router - same channel is used for each type but struct
	// defines unidirectionality of channel.
	r := &GenericRouter{
		externalMsgChan: msgChan,
		internalMsgChan: msgChan,
		externalRtChan:  rtChan,
		internalRtChan:  rtChan,
		externalRegChan: cmpChan,
		internalRegChan: cmpChan,
		rt:              rt,
		rc:              rc,
	}

	return r
}

// Send is an endpoint used for external go routines to send messages into the
// router
func (r *GenericRouter) Send(m interface{}) error {

	select {
	case r.externalMsgChan <- m:
		return nil
	default:
		return errors.New("Could not send message to router")
	}

}

// internal send method for routing messages to correct destinations
func (r *GenericRouter) send(m interface{}) error {

	//TODO: lookup source of sending message in routing table:
	// handle source not being in routing table, router needs to log this
	// if source is in routing table lookup associated []component array, cycles
	// through destinations and call their Send() functions.

}

// RegisterComponent registers a Component object with the router. This then
// allows us to use this component in route definitions
func (r *GenericRouter) RegisterComponent(c Component) {

}
