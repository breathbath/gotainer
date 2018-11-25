package container

import (
	"fmt"
	"strings"
)

//Tree of dependency nodes
type Tree []Node

//Event about registration of a specific service
type Event struct {
	Name    string
	Service string
}

func (e Event) String() string {
	return fmt.Sprintf(
		"{Name: %s; Service: %s;}",
		e.Name,
		e.Service,
	)
}

//Services list of dependencies
type Services []string

func (ss Services) String() string {
	return "[" + strings.Join(ss, ";") + "]"
}

//Observer is a service which is interested other services under a certain event
type Observer struct {
	Event    string
	Name     string
	Callback interface{}
}

//ParametersProvider gives container parameters
type ParametersProvider interface {
	GetItems() map[string]interface{}
}

func (o Observer) String() string {
	return fmt.Sprintf(
		"{Name: %s; Event: %s;}",
		o.Name,
		o.Event,
	)
}

//Node of a dependency
type Node struct {
	ID            string
	Constr        Constructor
	NewFunc       interface{}
	ServiceNames  Services
	Ev            Event
	Ob            Observer
	Parameters    map[string]interface{}
	ParamProvider ParametersProvider
	GarbageFunc   GarbageCollectorFunc
}

func (n Node) String() string {
	return fmt.Sprintf(
		"Node: {ID: %s; ServiceNames: %s; Event: %s; Observer: %s}",
		n.ID,
		n.ServiceNames,
		n.Ev,
		n.Ob,
	)
}

//ServiceExists checks if the provided name was already registered for a service
func (t Tree) ServiceExists(serviceID string) bool {
	for _, node := range t {
		if node.ID == serviceID {
			return true
		}
	}
	return false
}
