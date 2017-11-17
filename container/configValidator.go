package container

import (
	"fmt"
)

type NodeProcessor func(node Node, serviceId string) error

func validateNode(node Node, serviceId string, errCollection *[]error) {
	if node.NewFunc != nil {
		validateNewFunc(node, serviceId, errCollection)
	}

	if node.Constr != nil {
		validateConstrFunc(node, serviceId, errCollection)
	}

	if len(node.Ss) > 0 && node.NewFunc == nil {
		addErrorToCollection(errCollection, "Services list should be defined with a non empty new func", serviceId)
	}

	if node.Ob.Name != "" || node.Ob.Callback != nil || node.Ob.Event != "" {
		validateObserverDefinition(node, serviceId, errCollection)
	}

	if node.Ev.Name != "" || node.Ev.Service != "" {
		validateEventDefinition(node, serviceId, errCollection)
	}
}

func addErrorIfNeeded(errors *[]error, err error) {
	errs := *errors
	if err != nil {
		errs = append(errs, err)
	}
}

func ValidateConfig(tree Tree, nodeProcessor NodeProcessor) []error {
	errors := []error{}
	for serviceId, node := range tree {
		validateNode(node, serviceId, &errors)
		if err != nil {
			errors = append(errors, err)
		} else {
			nodeProcessor(node, serviceId)
		}
	}

	return errors
}

func validateNewFunc(node Node, serviceId string, errCollection *[]error) {
	assertConstructorIsEmpty(node, serviceId, errCollection)
	assertEventIsEmpty(node, serviceId, errCollection)
	assertObserverIsEmpty(node, serviceId, errCollection)
}

func validateConstrFunc(node Node, serviceId string, errCollection *[]error) {
	assertNewIsEmpty(node, serviceId, errCollection)
	assertEventIsEmpty(node, serviceId, errCollection)
	assertObserverIsEmpty(node, serviceId, errCollection)
}

func validateObserverDefinition(node Node, serviceId string, errCollection *[]error) {
	if node.Ob.Name == ""{
		addErrorToCollection(errCollection, "Observer name is required", serviceId)
	}
	if node.Ob.Callback == nil{
		addErrorToCollection(errCollection, "Observer callback is required", serviceId)
	}
	if node.Ob.Event == ""{
		addErrorToCollection(errCollection, "Observer event is required", serviceId)
	}
	assertNewIsEmpty(node, serviceId, errCollection)
	assertEventIsEmpty(node, serviceId, errCollection)
	assertConstructorIsEmpty(node, serviceId, errCollection)
}

func validateEventDefinition(node Node, serviceId string, errCollection *[]error) {
	if node.Ev.Name == ""{
		addErrorToCollection(errCollection, "Event name is required", serviceId)
	}
	if node.Ev.Service == ""{
		addErrorToCollection(errCollection, "Event service is required", serviceId)
	}

	assertNewIsEmpty(node, serviceId, errCollection)
	assertObserverIsEmpty(node, serviceId, errCollection)
	assertConstructorIsEmpty(node, serviceId, errCollection)
}

func assertNewIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.NewFunc != nil {
		addErrorToCollection(errCollection, "Unexpected new func declaration", serviceId)
	}
}

func assertConstructorIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.Constr != nil {
		addErrorToCollection(errCollection, "Unexpected constructor declaration", serviceId)
	}
}

func assertEventIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.Ev.Name != "" {
		addErrorToCollection(errCollection, "Unexpected event declaration", serviceId)
	}
}

func assertObserverIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.Ob.Name != "" {
		addErrorToCollection(errCollection, "Unexpected observer declaration", serviceId)
	}
}

func assertServicesIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if len(node.Ss) == 0 {
		addErrorToCollection(errCollection, "Unexpected services declaration", serviceId)
	}
}

func addErrorToCollection (errCollection *[]error, errorText, serviceId string) {
	*errCollection = append(*errCollection, fmt.Errorf(errorText + ", see service '%s'", serviceId))
}