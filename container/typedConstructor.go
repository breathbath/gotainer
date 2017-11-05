package container

import "reflect"

func convertTypedToUntypedConstructor(container Container, constructorFunc interface{}, constructorArgumentNames []string) ArgumentsConstructor {
	reflectedConstructorFunc := reflect.ValueOf(constructorFunc)

	assertFunctionDeclaredAsConstructor(reflectedConstructorFunc, constructorArgumentNames)
	assertFunctionReturnValues(reflectedConstructorFunc)

	argumentsToCallConstructorFunc := getValidFunctionArguments(reflectedConstructorFunc, constructorArgumentNames, container)

	return func(c Container) (interface{}, error) {
		values := reflectedConstructorFunc.Call(argumentsToCallConstructorFunc)
		if reflectedConstructorFunc.Type().NumOut() == 2 {
			if isErrorType(reflectedConstructorFunc.Type().Out(0)) {
				err := values[0].Interface().(error)
				obj := values[1].Interface()
				return obj, err
			}
			obj := values[0].Interface()
			err := values[1].Interface().(error)
			return obj, err
		}
		obj := values[0].Interface()
		return obj, nil
	}
}

func getValidFunctionArguments(reflectedConstructorFunc reflect.Value, constructorArgumentNames []string, container Container) []reflect.Value {
	constructorInputCount := reflectedConstructorFunc.Type().NumIn()
	argumentsToCallConstructorFunc := make([]reflect.Value, constructorInputCount)

	for i := 0; i < constructorInputCount; i++ {
		reflectedConstructorArgument := reflectedConstructorFunc.Type().In(i)

		correspondingServiceName := constructorArgumentNames[i]
		correspondingServiceFromContainer := container.GetService(correspondingServiceName)
		reflectedCorrespondingServiceFromContainer := reflect.ValueOf(correspondingServiceFromContainer)

		assertConstructorArgumentsAreCompatible(reflectedConstructorArgument, reflectedCorrespondingServiceFromContainer, correspondingServiceName)
		argumentsToCallConstructorFunc[i] = reflectedCorrespondingServiceFromContainer
	}

	return argumentsToCallConstructorFunc
}
