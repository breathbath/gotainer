package container

import (
	"fmt"
	"strings"
)

type Tree []Node

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

type Services []string

func (ss Services) String() string {
	return "[" + strings.Join(ss, ";") + "]"
}

type Observer struct {
	Event    string
	Name     string
	Callback interface{}
}

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

type Node struct {
	Id            string
	Constr        Constructor
	NewFunc       interface{}
	ServiceNames  Services
	Ev            Event
	Ob            Observer
	Parameters    map[string]interface{}
	ParamProvider ParametersProvider
}

func (n Node) String() string {
	return fmt.Sprintf(
		"Node: {Id: %s; ServiceNames: %s; Event: %s; Observer: %s}",
		n.Id,
		n.ServiceNames,
		n.Ev,
		n.Ob,
	)
}

func (t Tree) ServiceExists(serviceId string) bool {
	for _, node := range t {
		if node.Id == serviceId {
			return true
		}
	}
	return false
}
