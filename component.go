package msgrouter

import "errors"

type ComponentID int

type Component interface {
	Send() error
	GetID() error
}

type GenericComponent struct {
	ID int
}

func (c *GenericComponent) Send() error {
	return errors.New("Place Holder")
}

func (c *GenericComponent) GetID() error {
	return errors.New("Place Holder")
}
