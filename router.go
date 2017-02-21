package msgrouter

import "errors"

type Router interface {
	Send() error
	RegisterComponent(*Component) error
	UnregisterComponent(*Component) error
	AddRoute() error
	RemoveRoute() error
	ListRoutes() error
}

type GenericRouter struct {
	inChan               <-chan interface{}
	rt                   *routingTable
	registeredComponents []*Component
}

func NewGenericRouter() *GenericRouter {
	// Create inChan
	c := make(chan interface{}, 15000)

	// Create empty routing table
	var rt *routingTable

	// Instantiate GenericRouter
	r := &GenericRouter{
		inChan: c,
		rt:     rt,
	}

	return r

}

func (r *GenericRouter) Send() error {
	return errors.New("Place Holder")
}

func (r *GenericRouter) RegisterComponent(c *Component) error {
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
	return errors.New("Place Holder")
}

func (r *GenericRouter) RemoveRoute() error {
	return errors.New("Place Holder")
}

func (r *GenericRouter) ListRoutes() error {
	return errors.New("Place Holder")
}
