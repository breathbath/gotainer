package container

import (
	"fmt"
	"reflect"
)

func validateNode(node Node, serviceId string, errCollection *[]error, tree Tree) {
	if node.NewFunc != nil {
		validateNewFunc(node, serviceId, errCollection)
		return
	}

	if node.Constr != nil {
		validateConstrFunc(node, serviceId, errCollection)
		return
	}

	if len(node.ServiceNames) > 0 && node.NewFunc == nil {
		registerNewErrorInCollection(errCollection, "Services list should be defined with a non empty new func", serviceId)
		return
	}

	if node.Ob.Name != "" || node.Ob.Callback != nil || node.Ob.Event != "" {
		validateObserverDefinition(node, serviceId, errCollection)
		return
	}

	if node.Ev.Name != "" || node.Ev.Service != "" {
		validateEventDefinition(node, serviceId, errCollection, tree)
		return
	}
}

//ValidateConfig validates a tree of config options
func ValidateConfig(tree Tree) {
	errors := []error{}
	for serviceId, node := range tree {
		validateNode(node, serviceId, &errors, tree)
	}

	panicIfErrors(errors)
}

func validateNewFunc(node Node, serviceId string, errCollection *[]error) {
	assertConstructorIsEmpty(node, serviceId, errCollection)
	assertEventIsEmpty(node, serviceId, errCollection)
	assertObserverIsEmpty(node, serviceId, errCollection)

	var err error
	reflectedNewMethod := reflect.ValueOf(node.NewFunc)

	err = assertFunctionDeclaration(reflectedNewMethod, len(node.ServiceNames), serviceId)
	addErrorToCollection(errCollection, err)

	err = validateConstructorReturnValues(reflectedNewMethod, serviceId)
	addErrorToCollection(errCollection, err)
}

func validateConstrFunc(node Node, serviceId string, errCollection *[]error) {
	assertNewIsEmpty(node, serviceId, errCollection)
	assertEventIsEmpty(node, serviceId, errCollection)
	assertObserverIsEmpty(node, serviceId, errCollection)
}

func validateObserverDefinition(node Node, serviceId string, errCollection *[]error) {
	if node.Ob.Name == "" {
		registerNewErrorInCollection(errCollection, "Observer name is required", serviceId)
	}

	if node.Ob.Callback == nil {
		registerNewErrorInCollection(errCollection, "Observer callback is required", serviceId)
	}

	if node.Ob.Event == "" {
		registerNewErrorInCollection(errCollection, "Observer event is required", serviceId)
	}

	assertNewIsEmpty(node, serviceId, errCollection)
	assertEventIsEmpty(node, serviceId, errCollection)
	assertConstructorIsEmpty(node, serviceId, errCollection)

	err := assertFunctionDeclaration(reflect.ValueOf(node.Ob.Callback), 2, serviceId)
	addErrorToCollection(errCollection, err)
}

func validateEventDefinition(node Node, serviceId string, errCollection *[]error, tree Tree) {
	if node.Ev.Name == "" {
		registerNewErrorInCollection(errCollection, "Event name is required", serviceId)
	}
	if node.Ev.Service == "" {
		registerNewErrorInCollection(errCollection, "Event service is required", serviceId)
	}

	assertNewIsEmpty(node, serviceId, errCollection)
	assertObserverIsEmpty(node, serviceId, errCollection)
	assertConstructorIsEmpty(node, serviceId, errCollection)

	if node.Ev.Service != "" {
		assertServiceIsDeclared(node.Ev.Service, "event "+node.Ev.Name, tree, errCollection)
	}
}

func assertNewIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.NewFunc != nil {
		registerNewErrorInCollection(errCollection, "Unexpected new func declaration", serviceId)
	}
}

func assertConstructorIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.Constr != nil {
		registerNewErrorInCollection(errCollection, "Unexpected constructor declaration", serviceId)
	}
}

func assertEventIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.Ev.Name != "" {
		registerNewErrorInCollection(errCollection, "Unexpected event declaration", serviceId)
	}
}

func assertObserverIsEmpty(node Node, serviceId string, errCollection *[]error) {
	if node.Ob.Name != "" {
		registerNewErrorInCollection(errCollection, "Unexpected observer declaration", serviceId)
	}
}

func registerNewErrorInCollection(errCollection *[]error, errorText, serviceId string) {
	*errCollection = append(*errCollection, fmt.Errorf(errorText+", see service '%s'", serviceId))
}

func addErrorToCollection(errCollection *[]error, err error) {
	if err != nil {
		*errCollection = append(*errCollection, err)
	}
}

func assertServiceIsDeclared(serviceName, declarationPlace string, tree Tree, errCollection *[]error) {
	_, exists := tree[serviceName]
	if !exists {
		err := fmt.Errorf("Unknown service declaration '%s' in '%s'", serviceName, declarationPlace)
		addErrorToCollection(errCollection, err)
	}
}
