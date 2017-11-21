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
		return errors.New("must pass a pointer, not a value")
	}
	//if destinationPointerValue.IsNil() {
	//	return errors.New("nil pointer passed to destination")
	//}

	reflectedCreatedDependency := reflect.ValueOf(createdDependency)

	destinationValue := reflect.Indirect(destinationPointerValue)
	destinationValueType := destinationValue.Type()

	if reflectedCreatedDependency.Kind() == reflect.Ptr && sourceCanBeCopiedToDestination(reflectedCreatedDependency, destinationPointerValue) {
		reflectedCreatedDependencyIndirected := reflect.Indirect(reflectedCreatedDependency)
		convertedDependencyToDestinationValue := reflectedCreatedDependencyIndirected.Convert(destinationValueType)
		destinationValue.Set(convertedDependencyToDestinationValue)
		return nil
	}

	if sourceCanBeCopiedToDestination(reflectedCreatedDependency, destinationValue) {
		destinationValue.Set(reflectedCreatedDependency.Convert(destinationValueType))
		return nil
	}

	errStr := fmt.Sprintf(
		"Cannot convert created value of type '%s' to expected destination value '%s' for createdDependency declaration %s",
		reflectedCreatedDependency.Type().Name(),
		destinationValue.Type().Name(),
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
