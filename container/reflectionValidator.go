package container

import (
	"errors"
	"fmt"
	"reflect"
)

func isFunction(reflectedConstructorFunc reflect.Value) bool {
	return reflectedConstructorFunc.Kind() == reflect.Func
}

func hasArgumentsCount(reflectedConstructorFunc reflect.Value, expectedArgNumbersCount int) (bool, int) {
	if reflectedConstructorFunc.Type().IsVariadic() {
		return true, expectedArgNumbersCount
	}
	argsInputCount := reflectedConstructorFunc.Type().NumIn()
	return argsInputCount == expectedArgNumbersCount, argsInputCount
}

func assertFunctionDeclaration(
	reflectedConstructorFunc reflect.Value,
	expectedArgumentsCount int,
	serviceId string,
) error {
	if !isFunction(reflectedConstructorFunc) {
		errName := fmt.Sprintf(
			"A function is expected rather than '%s' [check '%s' service]",
			reflectedConstructorFunc.Kind(),
			serviceId,
		)
		return errors.New(errName)
	}

	hasArgCount, argsCount := hasArgumentsCount(reflectedConstructorFunc, expectedArgumentsCount)

	if !hasArgCount {
		errName := fmt.Sprintf(
			"The function requires %d arguments, but %d arguments are provided [check '%s' service]",
			argsCount,
			expectedArgumentsCount,
			serviceId,
		)
		return errors.New(errName)
	}

	return nil
}

func validateConstructorReturnValues(reflectedConstructorFunc reflect.Value, serviceId string) error {
	if reflectedConstructorFunc.Kind() != reflect.Func {
		return nil
	}

	constructorReturnsCount := reflectedConstructorFunc.Type().NumOut()

	if constructorReturnsCount > 2 || constructorReturnsCount < 1 {
		errName := fmt.Sprintf(
			"Constr function should return 1 or 2 values, but %d values are returned [check '%s' service]",
			constructorReturnsCount,
			serviceId,
		)
		return errors.New(errName)
	}

	var firstReflectedReturnValue, secondReflectedReturnValue reflect.Type

	if constructorReturnsCount == 2 {
		firstReflectedReturnValue = reflectedConstructorFunc.Type().Out(0)
		secondReflectedReturnValue = reflectedConstructorFunc.Type().Out(1)

		if !isErrorType(secondReflectedReturnValue) && !isErrorType(firstReflectedReturnValue) {
			return fmt.Errorf(
				"Constr function with 2 returned values should return at least one error interface [check '%s' service]",
				serviceId,
			)
		}
	}

	return nil
}

func assertConstructorArgumentsAreCompatible(
	reflectedConstructorArgument reflect.Type,
	reflectedContainerDependency reflect.Value,
	dependencyName, serviceId string,
) error {
	if reflectedConstructorArgument.Kind() == reflect.Interface && reflectedContainerDependency.Type().Implements(reflectedConstructorArgument) {
		return nil
	}

	realValueIsNil := reflectedContainerDependency.IsValid()
	if (reflectedConstructorArgument.Kind() == reflect.Interface || reflectedConstructorArgument.Kind() == reflect.Ptr) && !realValueIsNil {
		return nil
	}
	if reflectedConstructorArgument.Kind() != reflectedContainerDependency.Kind() ||
		!reflectedConstructorArgument.ConvertibleTo(reflectedContainerDependency.Type()) {
		errName := fmt.Sprintf(
			"Cannot use the provided dependency '%s' of type '%s' as '%s' in the Constr function call [check '%s' service]",
			dependencyName,
			reflectedContainerDependency.Type(),
			reflectedConstructorArgument,
			serviceId,
		)
		return errors.New(errName)
	}

	return nil
}
