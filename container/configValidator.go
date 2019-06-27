package container

import (
	"errors"
	"fmt"
	"reflect"
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

	assertConstructorOrNewFunctionAreDeclared(node, errCollection)

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
	}
}

//ValidateConfig validates a tree of config options and panics if something is wrong
func ValidateConfig(tree Tree) {
	err := ValidateConfigSecure(tree)

	panicIfError(err)
}

//ValidateConfigSecure validates a tree of config options and returns error if something is wrong
func ValidateConfigSecure(tree Tree) error {
	errs := []error{}
	for _, node := range tree {
		validateNode(node, &errs, tree)
	}

	return mergeErrors(errs)
}

func validateNewFunc(node Node, errCollection *[]error) {
	assertConstructorIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertObserverIsEmpty(node, errCollection)

	var err error
	reflectedNewMethod := reflect.ValueOf(node.NewFunc)

	err = assertFunctionDeclaration(reflectedNewMethod, len(node.ServiceNames), node.String())
	addErrorToCollection(errCollection, err)

	err = validateConstructorReturnValues(reflectedNewMethod, node.ID)
	addErrorToCollection(errCollection, err)
	assertServiceIDIsNotEmpty(node, errCollection, "The new function should be provided with a service id, see '%s'")
}

func validateConstrFunc(node Node, errCollection *[]error) {
	assertNewIsEmpty(node, errCollection)
	assertEventIsEmpty(node, errCollection)
	assertObserverIsEmpty(node, errCollection)
	assertServiceIDIsNotEmpty(node, errCollection, "The constructor function should be provided with a non empty service id, see '%s'")
}

func validateObserverDefinition(node Node, errCollection *[]error) {
	if node.Ob.Name == "" {
		registerNewErrorInCollection(errCollection, "Observer name is required [check '%s' service]", node)
	}

	if node.Ob.Callback == nil {
		registerNewErrorInCollection(errCollection, "Observer callback is required [check '%s' service]", node)
	}

	if node.Ob.Event == "" {
		registerNewErrorInCollection(errCollection, "Observer event is required [check '%s' service]", node)
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

	eventServiceIsFound, eventObserverIsFound := false, false
	if node.Ev.Service == "" {
		eventServiceIsFound = true //no need to look for empty service we skip the search function like this
	}
	if node.Ev.Name == "" {
		eventObserverIsFound = true //no need to look for empty event name we skip the search function like this
	}
	for _, searcheableNode := range tree {
		if eventServiceIsFound && eventObserverIsFound {
			break
		}

		if searcheableNode.ID == node.Ev.Service {
			eventServiceIsFound = true
		}

		if !searcheableNode.Ob.IsEmpty() && searcheableNode.Ob.Event == node.Ev.Name {
			eventObserverIsFound = true
		}
	}
	if !eventServiceIsFound {
		addErrorToCollection(
			errCollection,
			fmt.Errorf("Unknown service declaration '%s' in '%s'", node.Ev.Service, "event "+node.Ev.Name),
		)
	}

	if !eventObserverIsFound {
		addErrorToCollection(
			errCollection,
			fmt.Errorf("No observer is declared for the event '%s'", node.Ev.Name),
		)
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

func assertServiceIDIsNotEmpty(node Node, errCollection *[]error, errorFormat string) {
	if node.ID == "" {
		registerNewErrorInCollection(errCollection, errorFormat, node)
	}
}

func assertConstructorOrNewFunctionAreDeclared(node Node, errCollection *[]error) {
	if len(node.Parameters) == 0 && node.Ob.IsEmpty() && node.Ev.IsEmpty() && node.Constr == nil && node.NewFunc == nil {
		err := fmt.Errorf("A new or constructor function are expected but none was declared [check '%s' service]", node.ID)
		if node.ID == "" {
			err = fmt.Errorf("A new or constructor function are expected but none was declared see '%s'", node)
		}

		addErrorToCollection(errCollection, err)
	}
}
