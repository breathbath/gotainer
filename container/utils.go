package container

import (
	"errors"
	"fmt"
	"reflect"
)

//copySourceVariableToDestinationVariable copies a dependency fetched from the container
//to the pointer reference provided as dest in Scan or ScanNonCached of the container
func copySourceVariableToDestinationVariable(createdDependency interface{}, destination interface{}, dependencyName string) error {
	destinationPointerValue := reflect.ValueOf(destination)
	if destinationPointerValue.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"Please provide a pointer variable rather than a value [check '%s' service]",
			dependencyName,
		)
	}
	if destinationPointerValue.IsNil() {
		return fmt.Errorf(
			"Please provide an initialsed variable rather than a non-initialised pointer variable [check '%s' service]",
			dependencyName,
		)
	}

	reflectedCreatedDependency := reflect.ValueOf(createdDependency)

	destinationValue := reflect.Indirect(destinationPointerValue)

	if reflectedCreatedDependency.Kind() == reflect.Ptr && sourceCanBeCopiedToDestination(reflectedCreatedDependency, destinationPointerValue) {
		reflectedCreatedDependencyIndirected := reflect.Indirect(reflectedCreatedDependency)
		destinationValue.Set(reflectedCreatedDependencyIndirected.Convert(destinationValue.Type()))
		return nil
	}

	if sourceCanBeCopiedToDestination(reflectedCreatedDependency, destinationValue) {
		destinationValue.Set(reflectedCreatedDependency.Convert(destinationValue.Type()))
		return nil
	}

	errStr := fmt.Sprintf(
		"Cannot convert created value of type '%s' to expected destination value '%s' for createdDependency declaration %s [check '%s' service]",
		reflectedCreatedDependency.Type().Name(),
		destinationValue.Type().Name(),
		dependencyName,
		dependencyName,
	)

	return errors.New(errStr)
}

//sourceCanBeCopiedToDestination validates if we can copy the value received from the container to the defined dest pointer
func sourceCanBeCopiedToDestination(sourceValue, destinationValue reflect.Value) bool {
	destinationValueKind := destinationValue.Kind()
	sourceValueKind := sourceValue.Kind()
	isConvertable := sourceValue.Type().ConvertibleTo(destinationValue.Type())

	return destinationValueKind == sourceValueKind && isConvertable
}

func isErrorType(reflectedType reflect.Type) bool {
	expectedErrorInterface := reflect.TypeOf((*error)(nil)).Elem()
	return reflectedType.Implements(expectedErrorInterface)
}
