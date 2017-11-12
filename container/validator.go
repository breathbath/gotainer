package container

import (
	"errors"
	"fmt"
	"reflect"
)

func assertIsFunction(reflectedConstructorFunc reflect.Value) {
	if reflectedConstructorFunc.Kind() != reflect.Func {
		errName := fmt.Sprintf(
			"Destination object should be a constructor function rather than %s",
			reflectedConstructorFunc.Kind(),
		)
		panic(errors.New(errName))
	}
}

func assertArgumentsCount(reflectedConstructorFunc reflect.Value, expectedArgNumbersCount int) {
	argsInputCount := reflectedConstructorFunc.Type().NumIn()
	if argsInputCount != expectedArgNumbersCount {
		errName := fmt.Sprintf(
			"The function requires %d dependencies, but %d arguments are provided",
			argsInputCount,
			expectedArgNumbersCount,
		)
		panic(errors.New(errName))
	}
}

func assertFunctionDeclaredAsConstructor(reflectedConstructorFunc reflect.Value, constructorArgumentNames []string) {
	assertIsFunction(reflectedConstructorFunc)
	assertArgumentsCount(reflectedConstructorFunc, len(constructorArgumentNames))
}

func assertFunctionReturnValues(reflectedConstructorFunc reflect.Value) {
	constructorReturnsCount := reflectedConstructorFunc.Type().NumOut()

	if constructorReturnsCount > 2 || constructorReturnsCount < 1 {
		errName := fmt.Sprintf(
			"constructor function should return 1 or 2 values, but %d values are returned",
			constructorReturnsCount,
		)
		panic(errors.New(errName))
	}

	var firstReflectedReturnValue, secondReflectedReturnValue reflect.Type

	if constructorReturnsCount == 2 {
		firstReflectedReturnValue = reflectedConstructorFunc.Type().Out(0)
		secondReflectedReturnValue = reflectedConstructorFunc.Type().Out(1)

		if !isErrorType(secondReflectedReturnValue) && !isErrorType(firstReflectedReturnValue) {
			panic(errors.New("constructor function with 2 returned values should return at least one error interface"))
		}
	}
}

func assertConstructorArgumentsAreCompatible(
	reflectedConstructorArgument reflect.Type,
	reflectedContainerDependency reflect.Value,
	dependencyName string) {
	if reflectedConstructorArgument.Kind() == reflect.Interface && reflectedContainerDependency.Type().Implements(reflectedConstructorArgument) {
		return
	}
	if reflectedConstructorArgument.Kind() != reflectedContainerDependency.Kind() ||
		!reflectedConstructorArgument.ConvertibleTo(reflectedContainerDependency.Type()) {
		errName := fmt.Sprintf(
			"Cannot use the provided dependency '%s' of type '%s' as '%s' in the constructor function call",
			dependencyName,
			reflectedContainerDependency.Type(),
			reflectedConstructorArgument,
		)
		panic(errors.New(errName))
	}
}
