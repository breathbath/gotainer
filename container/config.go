package container

type Tree map[string]Node

type Event struct {
	Name    string
	Service string
}

type Services []string

type Observer struct {
	Event    string
	Name     string
	Callback interface{}
}

type Node struct {
	Constr  Constructor
	NewFunc interface{}
	Ss      Services
	Ev      Event
	Ob      Observer
}
