package container

import (
	"reflect"
	"fmt"
	"errors"
)

func assertFunctionDeclaredAsConstructor(reflectedConstructorFunc reflect.Value, constructorArgumentNames []string) {
	if reflectedConstructorFunc.Kind() != reflect.Func {
		errName := fmt.Sprintf(
			"Destination object should be a constructor function rather than %s",
			reflectedConstructorFunc.Kind(),
		)
		panic(errors.New(errName))
	}

	constructorInputCount := reflectedConstructorFunc.Type().NumIn()

	if constructorInputCount != len(constructorArgumentNames) {
		errName := fmt.Sprintf(
			"The constructor function requires %d dependencies, but %d argument names are provided",
			constructorInputCount,
			len(constructorArgumentNames),
		)
		panic(errors.New(errName))
	}
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

func assertConstructorArgumentsAreCompatible(reflectedConstructorArgument reflect.Type, reflectedCorrespondingServiceFromContainer reflect.Value, correspondingServiceName string) {
	if reflectedConstructorArgument.Kind() == reflect.Interface && reflectedCorrespondingServiceFromContainer.Type().Implements(reflectedConstructorArgument) {
		return
	}
	if reflectedConstructorArgument.Kind() != reflectedCorrespondingServiceFromContainer.Kind() ||
		!reflectedConstructorArgument.ConvertibleTo(reflectedCorrespondingServiceFromContainer.Type()) {
		errName := fmt.Sprintf(
			"Cannot use the provided service '%s' of type '%s' as '%s' in the constructor function call",
			correspondingServiceName,
			reflectedCorrespondingServiceFromContainer.Type(),
			reflectedConstructorArgument,
		)
		panic(errors.New(errName))
	}
}
