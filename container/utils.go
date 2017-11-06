package container

import (
	"reflect"
	"errors"
	"fmt"
)

func copySurceVariableToDestinationVariable(createdDependency interface{}, destination interface{}, serviceName string) error {
	destinationPointerValue := reflect.ValueOf(destination)
	if destinationPointerValue.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value")
	}
	if destinationPointerValue.IsNil() {
		return errors.New("nil pointer passed to destination")
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
		"Cannot convert created value of type '%s' to expected destination value '%s' for createdDependency declaration %s",
		reflectedCreatedDependency.Type().Name(),
		destinationValue.Type().Name(),
		serviceName,
	)

	return errors.New(errStr)
}

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
