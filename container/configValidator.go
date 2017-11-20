package container

import (
	"fmt"
	"reflect"
	"errors"
)

func validateNode(node Node, errCollection *[]error, tree Tree) {
	if node.NewFunc != nil {
		validateNewFunc(node, errCollection)
		return
	}

	if node.Constr != nil {
		validateConstrFunc(node, errCollection)
		return
	}

	if len(node.ServiceNames) > 0 && node.NewFunc == nil {
		registerNewErrorInCollection(errCollection, "Services list should be defined with a non empty new func", node.Id)
		return
	}

	if node.Ob.Name != "" || node.Ob.Callback != nil || node.Ob.Event != "" {
		validateObserverDefinition(node, errCollection)
		return
	}

	if node.Ev.Name != "" || node.Ev.Service != "" {
		validateEventDefinition(node, errCollection, tree)
		return
	}
}

//ValidateConfig validates a tree of config options
func ValidateConfig(tree Tree) {
	errors := []error{}
	for _, node := range tree {
		validateNode(node, &errors, tree)
	}

	panicIfErrors(errors)
}

func validateNewFunc(node Node, errCollection *[]error) {
	assertConstructorIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertObserverIsEmpty(node, errCollection)

	var err error
	reflectedNewMethod := reflect.ValueOf(node.NewFunc)

	err = assertFunctionDeclaration(reflectedNewMethod, len(node.ServiceNames), node.Id)
	addErrorToCollection(errCollection, err)

	err = validateConstructorReturnValues(reflectedNewMethod, node.Id)
	addErrorToCollection(errCollection, err)
	assertServiceIdIsNotEmpty(node, errCollection, "The new function should be provided with a service id")
}

func validateConstrFunc(node Node, errCollection *[]error) {
	assertNewIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertObserverIsEmpty(node, errCollection)
	assertServiceIdIsNotEmpty(node, errCollection, "The constructor function should be provided with a non empty service id")
}

func validateObserverDefinition(node Node, errCollection *[]error) {
	observerText := fmt.Sprint(node.Ob)
	if node.Ob.Name == "" {
		registerNewErrorInCollection(errCollection, "Observer name is required", observerText)
	}

	if node.Ob.Callback == nil {
		registerNewErrorInCollection(errCollection, "Observer callback is required", observerText)
	}

	if node.Ob.Event == "" {
		registerNewErrorInCollection(errCollection, "Observer event is required", observerText)
	}

	assertNewIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertConstructorIsEmpty(node, errCollection)

	err := assertFunctionDeclaration(reflect.ValueOf(node.Ob.Callback), 2, observerText)
	addErrorToCollection(errCollection, err)
}

func validateEventDefinition(node Node, errCollection *[]error, tree Tree) {
	eventText := fmt.Sprint(node.Ev)
	if node.Ev.Name == "" {
		registerNewErrorInCollection(errCollection, "Event name is required", eventText)
	}
	if node.Ev.Service == "" {
		registerNewErrorInCollection(errCollection, "Event service is required", eventText)
	}

	assertNewIsEmpty(node, errCollection)
	assertObserverIsEmpty(node, errCollection)
	assertConstructorIsEmpty(node, errCollection)

	if node.Ev.Service != "" {
		assertServiceIsDeclared(node.Ev.Service, "event "+node.Ev.Name, tree, errCollection)
	}
}

func assertNewIsEmpty(node Node, errCollection *[]error) {
	if node.NewFunc != nil {
		registerNewErrorInCollection(errCollection, "Unexpected new func declaration", node.Id)
	}
}

func assertConstructorIsEmpty(node Node, errCollection *[]error) {
	if node.Constr != nil {
		registerNewErrorInCollection(errCollection, "Unexpected constructor declaration", node.Id)
	}
}

func assertEventIsEmpty(node Node, errCollection *[]error) {
	if node.Ev.Name != "" {
		registerNewErrorInCollection(errCollection, "Unexpected event declaration", node.Id)
	}
}

func assertObserverIsEmpty(node Node, errCollection *[]error) {
	if node.Ob.Name != "" {
		registerNewErrorInCollection(errCollection, "Unexpected observer declaration", node.Id)
	}
}

func registerNewErrorInCollection(errCollection *[]error, errorText, itemId string) {
	if itemId != "" {
		errorText = errorText + ", see '" + itemId + "'"
	}
	*errCollection = append(*errCollection, errors.New(errorText))
}

func addErrorToCollection(errCollection *[]error, err error) {
	if err != nil {
		*errCollection = append(*errCollection, err)
	}
}

func assertServiceIsDeclared(serviceName, declarationPlace string, tree Tree, errCollection *[]error) {
	if !tree.ServiceExists(serviceName) {
		err := fmt.Errorf("Unknown service declaration '%s' in '%s'", serviceName, declarationPlace)
		addErrorToCollection(errCollection, err)
	}
}

func assertServiceIdIsNotEmpty(node Node, errCollection *[]error, errorText string) {
	if node.Id == "" {
		nodeText := fmt.Sprint(node)
		registerNewErrorInCollection(errCollection, errorText, nodeText)
	}
}
