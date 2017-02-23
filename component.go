package msgrouter

type ComponentID int

type Component interface {
	Send() error
}
