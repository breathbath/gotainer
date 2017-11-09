package container

import "reflect"

func convertNewMethodToConstructor(container Container, newMethod interface{}, newMethodArgumentNames []string) Constructor {
	reflectedNewMethod := reflect.ValueOf(newMethod)

	assertFunctionDeclaredAsConstructor(reflectedNewMethod, newMethodArgumentNames)
	assertFunctionReturnValues(reflectedNewMethod)

	argumentsToCallConstructorFunc := getValidFunctionArguments(reflectedNewMethod, newMethodArgumentNames, container)

	return func(c Container) (interface{}, error) {
		values := reflectedNewMethod.Call(argumentsToCallConstructorFunc)
		if reflectedNewMethod.Type().NumOut() == 2 {
			if isErrorType(reflectedNewMethod.Type().Out(0)) {
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

func getValidFunctionArguments(reflectedNewMethod reflect.Value, newMethodArgumentNames []string, container Container) []reflect.Value {
	constructorInputCount := reflectedNewMethod.Type().NumIn()
	argumentsToCallNewMethod := make([]reflect.Value, constructorInputCount)

	for i := 0; i < constructorInputCount; i++ {
		reflectedNewMethodArgument := reflectedNewMethod.Type().In(i)

		dependencyName := newMethodArgumentNames[i]
		dependencyFromContainer := container.Get(dependencyName, true)
		reflectedDepdendencyFromContainer := reflect.ValueOf(dependencyFromContainer)

		assertConstructorArgumentsAreCompatible(reflectedNewMethodArgument, reflectedDepdendencyFromContainer, dependencyName)
		argumentsToCallNewMethod[i] = reflectedDepdendencyFromContainer
	}

	return argumentsToCallNewMethod
}
