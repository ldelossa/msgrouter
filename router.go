package msgrouter

import (
	"errors"
	"fmt"
)

// Router allows synchronization and routing decisions to be made between
// go routines. By having go routines send to a router we can create different
// messaging topologies. Operations on the router (route add, route remove,
// registering components, etc...) should block the receiving of messages from
// external go routines thus synchronizing updates without the need for locks
type Router interface {
	Send(msg *interface{}) error
	RegisterComponent(msgMsg) error
	UnregisterComponent(msgMsg) error
	AddRoute(msgRt) error
	RemoveRoute(msgRt) error
	// ListRoutes() (string, error)
	Consume()
}

// Operation constants to multiplex operations over channels

//REGISTER is a op code for msgReg. Tells router to use registerComponent handler
const REGISTER = 0

// UNREGISTER is an op code for msgReg. Tells router to use unregisterComponent
// handler
const UNREGISTER = 1

// ADDROUTE is an op code for msgRt. Tells router to use addRoute handler.
const ADDROUTE = 0

// REMOVEROUTE is an op code for msgRt. Tells router to use removeRoute handler.
const REMOVEROUTE = 1

// LISTROUTES is an op code for msgRt. Tells router to use listRoutes handler.
const LISTROUTES = 2

// ComponentID is an ID used to select registered components
type ComponentID UUID

// Map which correlates source component to one or more destination components
type routingTable map[ComponentID][]Component

// GenericRouter is an implementation of a router. External channels are for
// API access while internal channels are for consuming off of.
// TODO: make struct or interface to handle messages on each channel, removing
// interface{}
type GenericRouter struct {
	externalMsgChan chan<- msgMsg
	internalMsgChan <-chan msgMsg
	externalRtChan  chan<- msgRt
	internalRtChan  <-chan msgRt
	externalRegChan chan<- msgReg
	internalRegChan <-chan msgReg
	rt              routingTable
	rc              map[ComponentID]Component
}

// msg* structs are used to package messages that will be sent on the
// associated channel. External API
type msgMsg struct {
	src     ComponentID
	payload interface{}
}

type msgRt struct {
	op   int
	src  ComponentID
	dest ComponentID
}

type msgReg struct {
	c  Component
	op int
}

// NewGenericRouter is a constructor for a generic implementation of a Router
// Channels should be buffered so that sending go routines do not block while
// blocking operations occur on router
func NewGenericRouter(bufferSize int) *GenericRouter {

	// make channels
	msgChan := make(chan msgMsg, bufferSize)
	rtChan := make(chan msgRt, bufferSize)
	cmpChan := make(chan msgReg, bufferSize)

	// create routing table
	rt := routingTable{}

	// create registeredComponents map
	rc := make(map[ComponentID]Component)

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

// Consume is meant to be ran as a go routine. Consume will listen on all
// internal message channels and run the appropriate function handler based on the
// message received.
func (r *GenericRouter) Consume() {

	select {
	case m := <-r.internalMsgChan:
		go r.send(m)
	case m := <-r.internalRtChan:
		switch {
		case m.op == ADDROUTE:
			r.addRoute(m)
		case m.op == REMOVEROUTE:
			r.removeRoute(m)
		case m.op == LISTROUTES:
			fmt.Println()
		}
	case m := <-r.internalRegChan:
		switch {
		case m.op == UNREGISTER:
			r.unregisterComponent(m)
		case m.op == REGISTER:
			r.registerComponent(m)
		}
	}

}

// Send is a wrapper for external usage. Wrapping a send to the
// external message channel of our router.
func (r *GenericRouter) Send(m msgMsg) error {

	select {
	case r.externalMsgChan <- m:
		return nil
	default:
		return errors.New("Could not send message to router")
	}

}

func (r *GenericRouter) send(m msgMsg) {

	// Confirm src in msgMsg is in component array
	if _, ok := r.rc[m.src]; !ok {
		return
	}

	// Obtain routes
	routesArray, ok := r.rt[m.src]
	if !ok {
		return
	}

	// Send payload to each route.
	for _, comp := range routesArray {
		comp.Send(m.payload)
	}

}

// // internal send method for routing messages to correct destinations
// func (r *GenericRouter) send(m interface{}) error {
//
// 	//TODO: lookup source of sending message in routing table:
// 	// handle source not being in routing table, router needs to log this
// 	// if source is in routing table lookup associated []component array, cycles
// 	// through destinations and call their Send() functions.
//
// }

// func (r *GenericRouter) register(c Component) (ComponentID, error) {
// 	uuid, err := newUUID()
// 	if err != nil {
// 		return ComponentID(""), errors.New("Could not generate UUID")
// 	}
// 	c.SetID(uuid)
// 	r.rc[uuid] = c
// 	return uuid, nil
// }

// RegisterComponent is a wrapper for external usage. Wrapping a send to the
// external registration channel of our router.
func (r *GenericRouter) RegisterComponent(m msgReg) {
	// Tag on operation constant
	m.op = REGISTER
	// Send msgReg to external msgChan
	r.externalRegChan <- m

}

func (r GenericRouter) registerComponent(m msgReg) error {
	// Check to see if component already has ID
	id, err := m.c.GetID()
	if err == nil {

		// If component ID found, do lookup of ID in rc table.
		if comp, ok := r.rc[id]; ok {

			// Lookup of id succeeded, and component being registered matches lookup,
			// return hash, already registered.
			if comp == m.c {
				return nil
			}

		}
	}

	// This is a fallthrough. Didn't come in with ID or came in with ID but component
	// didn't match. Register and setID on component.
	uuid, err := newUUID()
	if err != nil {
		return errors.New("Could not generate UUID")
	}
	m.c.SetID(uuid)
	r.rc[uuid] = m.c
	return nil

}

// UnregisterComponent is a wrapper for external usage. Wrapping a send to the
// external unregistration channel of our router.
func (r *GenericRouter) UnregisterComponent(m msgReg) {
	// Tag on operation constant
	m.op = UNREGISTER
	// send msgReg to external msgChan
	r.externalRegChan <- m
}

// unregisterComponent searches the registeredComponent table for the hash
// that's in msgReg.Component. It will remove the component from the rc.
// This does not stop routing of registered components. Removing the route
// is necessary. TODO: Block removal of component if route exists for the
// component.
func (r GenericRouter) unregisterComponent(m msgReg) error {
	// Check to see if component has ID
	id, err := m.c.GetID()
	if err == nil {
		// If component has hash, look up hash in rc. If lookup succeeds, delete
		// the map entry
		if _, ok := r.rc[id]; ok {
			delete(r.rc, id)
			return nil
		}

	}
	return errors.New("Component not registered")
}

// AddRoute is a wrapper for external usage. Wrapping a send to the
// external route channel of our router.
func (r *GenericRouter) AddRoute(m msgRt) {
	// Tag on operation constant
	m.op = ADDROUTE
	// send msgRt to external msgChan
	r.externalRtChan <- m
}

// addRoute adds a component to an array of components. This array is hashed
// on the componetID, associating a component with it's routes. Only components
// registered by RegisterComponent are applicable for routes.
func (r *GenericRouter) addRoute(m msgRt) {

	// Confirm source is in registered components array
	if _, ok := r.rc[m.src]; !ok {
		return
	}
	if _, ok := r.rc[m.dest]; !ok {
		return
	}

	srcArray := r.rt[m.src]

	// Add destination component into source component's array. Lookup component
	// in registered component array
	srcArray = append(srcArray, r.rc[m.dest])

}

// RemoveRoute is a wrapper for external usage. Wrapping a send to the
// external route channel of our router.
func (r *GenericRouter) RemoveRoute(m msgRt) {
	// Tag on operation constant
	m.op = REMOVEROUTE
	// Send msgRt to external msgChan
	r.externalRtChan <- m
}

// removeRoute lookups a route's source, locates the given destination and
// removes this destination from the route's component array.
func (r *GenericRouter) removeRoute(m msgRt) {

	// Confirm source is in registered components array
	if _, ok := r.rc[m.src]; !ok {
		return
	}
	if _, ok := r.rc[m.dest]; !ok {
		return
	}

	// Lookup component array for source
	srcArray := r.rt[m.src]

	// Cycle through source array, remove destination component if found. Rrder
	// not important so just swap to last and return len - 1
	for i, c := range srcArray {
		if r.rc[m.dest] == c {
			srcArray[len(srcArray)-1], srcArray[i] = srcArray[i], srcArray[len(srcArray)-1]
			srcArray = srcArray[:len(srcArray)-1]
		}
	}

}
