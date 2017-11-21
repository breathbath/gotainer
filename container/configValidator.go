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
		registerNewErrorInCollection(errCollection, "Services list should be defined with a non empty new func, see '%s'", node)
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

	err = assertFunctionDeclaration(reflectedNewMethod, len(node.ServiceNames), node.String())
	addErrorToCollection(errCollection, err)

	err = validateConstructorReturnValues(reflectedNewMethod, node.Id)
	addErrorToCollection(errCollection, err)
	assertServiceIdIsNotEmpty(node, errCollection, "The new function should be provided with a service id, see '%s'")
}

func validateConstrFunc(node Node, errCollection *[]error) {
	assertNewIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertObserverIsEmpty(node, errCollection)
	assertServiceIdIsNotEmpty(node, errCollection, "The constructor function should be provided with a non empty service id, see '%s'")
}

func validateObserverDefinition(node Node, errCollection *[]error) {
	if node.Ob.Name == "" {
		registerNewErrorInCollection(errCollection, "Observer name is required, see '%s'", node)
	}

	if node.Ob.Callback == nil {
		registerNewErrorInCollection(errCollection, "Observer callback is required, see '%s'", node)
	}

	if node.Ob.Event == "" {
		registerNewErrorInCollection(errCollection, "Observer event is required, see '%s'", node)
	}

	assertNewIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertConstructorIsEmpty(node, errCollection)

	err := assertFunctionDeclaration(reflect.ValueOf(node.Ob.Callback), 2, node.String())
	addErrorToCollection(errCollection, err)
}

func validateEventDefinition(node Node, errCollection *[]error, tree Tree) {
	if node.Ev.Name == "" {
		registerNewErrorInCollection(errCollection, "Event name is required, see '%s'", node)
	}
	if node.Ev.Service == "" {
		registerNewErrorInCollection(errCollection, "Event service is required, see '%s'", node)
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
		registerNewErrorInCollection(errCollection, "Unexpected new func declaration, see '%s'", node)
	}
}

func assertConstructorIsEmpty(node Node, errCollection *[]error) {
	if node.Constr != nil {
		registerNewErrorInCollection(errCollection, "Unexpected constructor declaration, see '%s'", node)
	}
}

func assertEventIsEmpty(node Node, errCollection *[]error) {
	if node.Ev.Name != "" {
		registerNewErrorInCollection(errCollection, "Unexpected event declaration, see '%s'", node)
	}
}

func assertObserverIsEmpty(node Node, errCollection *[]error) {
	if node.Ob.Name != "" {
		registerNewErrorInCollection(errCollection, "Unexpected observer declaration, see '%s'", node)
	}
}

func registerNewErrorInCollection(errCollection *[]error, format string, context ...interface{}) {
	errorText := format
	if len(context) > 0 {
		errorText = fmt.Sprintf(format, context...)
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

func assertServiceIdIsNotEmpty(node Node, errCollection *[]error, errorFormat string) {
	if node.Id == "" {
		registerNewErrorInCollection(errCollection, errorFormat, node)
	}
}
