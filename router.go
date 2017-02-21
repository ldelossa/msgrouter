package msgrouter

import "errors"

type Router interface {
	Send(msg *interface{}) error
	RegisterComponent(c *Component) error
	UnregisterComponent(c *Component) error
	AddRoute(src *Component, dst *Component) error
	RemoveRoute(src *Component, dst *Component) error
	ListRoutes() (string, error)
}

type GenericRouter struct {
	inChan               chan<- interface{}
	rt                   *routingTable
	registeredComponents []*Component
}

func NewGenericRouter(bufferSize int) *GenericRouter {
	// Create inChan
	c := make(chan interface{}, bufferSize)

	// Create empty routing table
	var rt *routingTable

	// Instantiate GenericRouter
	r := &GenericRouter{
		inChan: c,
		rt:     rt,
	}

	return r

}

func (r *GenericRouter) Send(msg *interface{}) error {

	select {
	case r.inChan <- msg:
		return nil
	default:
		return errors.New("Could not send message to Router")
	}

}

func (r *GenericRouter) RegisterComponent(c *Component) error {
  // Search for component, if exists return this
	for _, ci := range r.registeredComponents {
		if ci == c 
	}


	// Add component to registeredComponents
	r.registeredComponents = append(r.registeredComponents, c)

	return errors.New("Place Holder")

}

func (r *GenericRouter) UnregisterComponent(c *Component) error {
	// Iterate through  registeredComponents and remove
	for i, ci := range r.registeredComponents {

		if c == ci {
			r.registeredComponents = append(r.registeredComponents[:i], r.registeredComponents[i+1:]...)
			return nil
		}
	}
	return errors.New("Component not registered")
}

func (r *GenericRouter) AddRoute() error {

}

func (r *GenericRouter) RemoveRoute() error {
	return errors.New("Place Holder")
}

func (r *GenericRouter) ListRoutes() error {
	return errors.New("Place Holder")
}
