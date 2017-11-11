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

/**
converts something like func(observer Observer, dependency Dependency) which is customObserverResolver to func(observer interface{}, dependency interface{}) which is dependencyNotifier
as customObserverResolver can be anything we need to make sure that function
 */
func convertCustomObserverResolverToDependencyNotifier(customObserverResolver interface{}, eventName, observerId string) dependencyNotifier {
	reflectedCustomObserverResolver := reflect.ValueOf(customObserverResolver)
	assertIsFunction(reflectedCustomObserverResolver)
	assertArgumentsCount(reflectedCustomObserverResolver, 2)
	return func(observer interface{}, dependency interface{}) {
		argumentsToCallCustomerObserverResolver := make([]reflect.Value, 2)

		reflectedObserver := reflect.ValueOf(observer)
		reflectedFirstResolverArgument := reflectedCustomObserverResolver.Type().In(0)
		assertConstructorArgumentsAreCompatible(reflectedFirstResolverArgument, reflectedObserver, observerId)
		argumentsToCallCustomerObserverResolver[0] = reflectedObserver

		reflectedDependency := reflect.ValueOf(dependency)
		reflectedSecondResolverArgument := reflectedCustomObserverResolver.Type().In(1)
		assertConstructorArgumentsAreCompatible(reflectedSecondResolverArgument, reflectedDependency, eventName)
		argumentsToCallCustomerObserverResolver[1] = reflectedDependency

		reflectedCustomObserverResolver.Call(argumentsToCallCustomerObserverResolver)
	}
}

func getValidFunctionArguments(reflectedNewMethod reflect.Value, newMethodArgumentNames []string, container Container) []reflect.Value {
	constructorInputCount := reflectedNewMethod.Type().NumIn()
	argumentsToCallNewMethod := make([]reflect.Value, constructorInputCount)

	for i := 0; i < constructorInputCount; i++ {
		reflectedNewMethodArgument := reflectedNewMethod.Type().In(i)

		dependencyName := newMethodArgumentNames[i]
		dependencyFromContainer := container.Get(dependencyName, true)
		reflectedDependencyFromContainer := reflect.ValueOf(dependencyFromContainer)

		assertConstructorArgumentsAreCompatible(reflectedNewMethodArgument, reflectedDependencyFromContainer, dependencyName)
		argumentsToCallNewMethod[i] = reflectedDependencyFromContainer
	}

	return argumentsToCallNewMethod
}
