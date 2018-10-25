package container

import "reflect"

//convertNewMethodToConstructor creates a Callback that will call a New method of a Service with the Config
//declared as newMethodArgumentNames.
//Suppose we have func NewServiceA(sb ServiceB, sc ServiceC) ServiceA, if you call
//convertNewMethodToConstructor(container, NewServiceA, "service_b", "service_c"), you will get a Callback that will:
//a) fetch "service_b" and "service_c" from the container
//b) validate if type of "service_b" and "service_c" is convertable to the NewServiceA arguments
//Constr) call NewServiceA with the results of container.Get("service_b") and container.Get("service_c")
func convertNewMethodToConstructor(
	container Container,
	newMethod interface{},
	newMethodArgumentNames []string,
	serviceId string,
) Constructor {

	reflectedNewMethod := reflect.ValueOf(newMethod)

	panicIfError(
		assertFunctionDeclaration(reflectedNewMethod, len(newMethodArgumentNames), serviceId),
	)

	panicIfError(
		validateConstructorReturnValues(reflectedNewMethod, serviceId),
	)

	return func(c Container) (interface{}, error) {
		argumentsToCallConstructorFunc, errors := getValidFunctionArguments(
			reflectedNewMethod,
			newMethodArgumentNames,
			container,
			serviceId,
		)
		panicIfErrors(errors)
		values := reflectedNewMethod.Call(argumentsToCallConstructorFunc)
		if reflectedNewMethod.Type().NumOut() == 2 {
			if isErrorType(reflectedNewMethod.Type().Out(0)) {
				return collectErrorAndResult(values[0], values[1])
			}
			return collectErrorAndResult(values[1], values[0])
		}
		return values[0].Interface(), nil
	}
}

func collectErrorAndResult(reflectedErrorValue, reflectedServiceValue reflect.Value) (interface{}, error) {
	err := getErrorOrNil(reflectedErrorValue)
	service := reflectedServiceValue.Interface()

	return service, err
}

func getErrorOrNil(value reflect.Value) error {
	var err error
	if value.IsNil() {
		err = nil
	} else {
		err = value.Interface().(error)
	}

	return err
}

//wrapCallbackToProvideDependencyToServiceIntoServiceNotificationCallback converts something like func(Observer Observer, dependency Dependency)
// which is customObserverResolver to func(Observer interface{}, dependency interface{}) which is serviceNotificationCallback
//as customObserverResolver can be anything we need to make sure that function
func wrapCallbackToProvideDependencyToServiceIntoServiceNotificationCallback(customObserverResolver interface{}, eventName, observerId string) serviceNotificationCallback {
	reflectedCustomObserverResolver := reflect.ValueOf(customObserverResolver)
	panicIfError(assertFunctionDeclaration(reflectedCustomObserverResolver, 2, observerId))

	//here we redirect a call to func(Observer interface{}, dependency interface{}) into
	// func(Observer Observer, dependency Dependency) which was given as customObserverResolver
	return func(observer interface{}, dependency interface{}) {
		argumentsToCallCustomerObserverResolver := make([]reflect.Value, 2)

		reflectedObserver := reflect.ValueOf(observer)
		reflectedFirstResolverArgument := reflectedCustomObserverResolver.Type().In(0)
		panicIfError(assertConstructorArgumentsAreCompatible(
			reflectedFirstResolverArgument, reflectedObserver, observerId, observerId),
		)
		argumentsToCallCustomerObserverResolver[0] = reflectedObserver

		reflectedDependency := reflect.ValueOf(dependency)
		reflectedSecondResolverArgument := reflectedCustomObserverResolver.Type().In(1)
		panicIfError(assertConstructorArgumentsAreCompatible(
			reflectedSecondResolverArgument, reflectedDependency, eventName, observerId),
		)
		argumentsToCallCustomerObserverResolver[1] = reflectedDependency

		reflectedCustomObserverResolver.Call(argumentsToCallCustomerObserverResolver)
	}
}

//getValidFunctionArguments fetches Config by ids defined in newMethodArgumentNames and validates if they are convertable
//to arguments of reflectedNewMethod which is a New method of a Service provided in the AddNewMethod of the container
func getValidFunctionArguments(reflectedNewMethod reflect.Value, newMethodArgumentNames []string, container Container, serviceId string) ([]reflect.Value, []error) {
	constructorInputCount := reflectedNewMethod.Type().NumIn()
	argumentsToCallNewMethod := make([]reflect.Value, constructorInputCount)

	var errors []error
	for i := 0; i < constructorInputCount; i++ {
		reflectedNewMethodArgument := reflectedNewMethod.Type().In(i)

		dependencyName := newMethodArgumentNames[i]
		dependencyFromContainer := container.Get(dependencyName, true)
		reflectedDependencyFromContainer := reflect.ValueOf(dependencyFromContainer)

		err := assertConstructorArgumentsAreCompatible(
			reflectedNewMethodArgument,
			reflectedDependencyFromContainer,
			dependencyName,
			serviceId,
		)
		if err != nil {
			errors = append(errors, err)
		} else {
			argumentsToCallNewMethod[i] = reflectedDependencyFromContainer
		}
	}

	return argumentsToCallNewMethod, errors
}
