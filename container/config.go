package container

type Tree []Node

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
	Id           string
	Constr       Constructor
	NewFunc      interface{}
	ServiceNames Services
	Ev           Event
	Ob           Observer
}

func (t Tree) ServiceExists(serviceId string) bool {
	for _, node := range t {
		if node.Id == serviceId {
			return true
		}
	}
	return false
}
