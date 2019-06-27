package container

import (
	"fmt"
	"reflect"
)

//validateDestinationPointerValue checks if destination value is not a pointer or a nil
func validateDestinationPointerValue(destinationPointerValue reflect.Value, dependencyName string) error {
	if destinationPointerValue.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"Please provide a pointer variable rather than a value [check '%s' service]",
			dependencyName,
		)
	}
	if destinationPointerValue.IsNil() {
		return fmt.Errorf(
			"Please provide an initialized variable rather than a non-initialised pointer variable [check '%s' service]",
			dependencyName,
		)
	}

	return nil
}

//copySourceVariableToDestinationVariable copies a dependency fetched from the container
//to the pointer reference provided as dest in Scan or ScanNonCached of the container
func copySourceVariableToDestinationVariable(
	createdDependency interface{},
	destination interface{},
	dependencyName string,
) error {
	destinationPointerValue := reflect.ValueOf(destination)
	err := validateDestinationPointerValue(destinationPointerValue, dependencyName)
	if err != nil {
		return err
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

	return fmt.Errorf(
		"Cannot convert created value of type '%s' to expected destination value '%s' for createdDependency declaration %s [check '%s' service]",
		reflectedCreatedDependency.Type().Name(),
		destinationValue.Type().Name(),
		dependencyName,
		dependencyName,
	)
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
