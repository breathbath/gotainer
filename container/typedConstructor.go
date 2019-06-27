package container

import (
	"reflect"
)

//convertNewMethodToNewFuncConstructor creates a Callback that will call a New method of a Service with the Config
//declared as newMethodArgumentNames.
//Suppose we have func NewServiceA(sb ServiceB, sc ServiceC) ServiceA, if you call
//convertNewMethodToNewFuncConstructor(container, NewServiceA, "service_b", "service_c"), you will get a Callback that will:
//a) fetch "service_b" and "service_c" from the container
//b) validate if type of "service_b" and "service_c" is convertable to the NewServiceA arguments
//Constr) call NewServiceA with the results of container.Get("service_b") and container.Get("service_c")
func convertNewMethodToNewFuncConstructor(
	container Container,
	newMethod interface{},
	newMethodArgumentNames []string,
	serviceId string,
) (NewFuncConstructor, error) {
	reflectedNewMethod := reflect.ValueOf(newMethod)

	err := assertFunctionDeclaration(reflectedNewMethod, len(newMethodArgumentNames), serviceId)
	if err != nil {
		return nil, err
	}

	err = validateConstructorReturnValues(reflectedNewMethod, serviceId)
	if err != nil {
		return nil, err
	}

	return func(c Container, isCached bool) (interface{}, error) {
		argumentsToCallConstructorFunc, err := getValidFunctionArguments(
			reflectedNewMethod,
			newMethodArgumentNames,
			container,
			serviceId,
			isCached,
		)
		if err != nil {
			return nil, err
		}

		values := reflectedNewMethod.Call(argumentsToCallConstructorFunc)
		if reflectedNewMethod.Type().NumOut() == 2 {
			if isErrorType(reflectedNewMethod.Type().Out(0)) {
				return collectErrorAndResult(values[0], values[1])
			}
			return collectErrorAndResult(values[1], values[0])
		}
		return values[0].Interface(), nil
	}, nil
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
func wrapCallbackToProvideDependencyToServiceIntoServiceNotificationCallback(
	customObserverResolver interface{},
	eventName,
	observerId string,
) (serviceNotificationCallback, error) {
	reflectedCustomObserverResolver := reflect.ValueOf(customObserverResolver)
	err := assertFunctionDeclaration(reflectedCustomObserverResolver, 2, observerId)
	if err != nil {
		return nil, err
	}

	//here we redirect a call to func(Observer interface{}, dependency interface{}) into
	// func(Observer Observer, dependency Dependency) which was given as customObserverResolver
	return func(observer interface{}, dependency interface{}) error {
		argumentsToCallCustomerObserverResolver := make([]reflect.Value, 2)

		reflectedObserver := reflect.ValueOf(observer)
		reflectedFirstResolverArgument := reflectedCustomObserverResolver.Type().In(0)
		err := assertCompatible(
			reflectedFirstResolverArgument,
			reflectedObserver.Type(),
			observerId,
			observerId,
		)
		if err != nil {
			return err
		}

		argumentsToCallCustomerObserverResolver[0] = reflectedObserver

		reflectedDependency := reflect.ValueOf(dependency)
		reflectedSecondResolverArgument := reflectedCustomObserverResolver.Type().In(1)
		err = assertCompatible(
			reflectedSecondResolverArgument,
			reflectedDependency.Type(),
			eventName,
			observerId,
		)
		if err != nil {
			return err
		}

		argumentsToCallCustomerObserverResolver[1] = reflectedDependency

		reflectedCustomObserverResolver.Call(argumentsToCallCustomerObserverResolver)
		return nil
	}, nil
}

//getValidFunctionArgumentsForVariadicFunc does the same as getValidFunctionArguments but processes a variadic
//constructor function
func getValidFunctionArgumentsForVariadicFunc(
	reflectedNewMethod reflect.Value,
	newMethodArgumentNames []string,
	container Container,
	serviceId string,
	isCached bool,
) ([]reflect.Value, error) {
	argumentsToCallNewMethod := []reflect.Value{}
	constructorInputCount := reflectedNewMethod.Type().NumIn()
	var i = 0
	var errors []error
	for _, dependencyName := range newMethodArgumentNames {
		i++
		dependencyFromContainer, err := container.GetSecure(dependencyName, isCached)
		if err != nil {
			return nil, err
		}

		reflectedDependencyFromContainer := reflect.ValueOf(dependencyFromContainer)

		var reflectedNewMethodArgument reflect.Type
		if i < constructorInputCount {
			reflectedNewMethodArgument = reflectedNewMethod.Type().In(i - 1)
		} else {
			reflectedVariadicArgumentCollection := reflectedNewMethod.Type().In(constructorInputCount - 1)
			reflectedNewMethodArgument = reflectedVariadicArgumentCollection.Elem()
		}

		reflectedDependencyFromContainer = replaceCompatibleNilDependency(
			reflectedNewMethodArgument,
			reflectedDependencyFromContainer,
			dependencyFromContainer,
		)

		err = assertCompatible(
			reflectedNewMethodArgument,
			reflectedDependencyFromContainer.Type(),
			dependencyName,
			serviceId,
		)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		argumentsToCallNewMethod = append(argumentsToCallNewMethod, reflectedDependencyFromContainer)
	}

	return argumentsToCallNewMethod, mergeErrors(errors)
}

//getValidFunctionArguments fetches Config by ids defined in newMethodArgumentNames and validates if they are convertable
//to arguments of reflectedNewMethod which is a New method of a Service provided in the AddNewMethod of the container
func getValidFunctionArguments(
	reflectedNewMethod reflect.Value,
	newMethodArgumentNames []string,
	container Container,
	serviceId string,
	isCached bool,
) ([]reflect.Value, error) {
	if reflectedNewMethod.Type().IsVariadic() {
		return getValidFunctionArgumentsForVariadicFunc(
			reflectedNewMethod,
			newMethodArgumentNames,
			container,
			serviceId,
			isCached,
		)
	}

	constructorInputCount := reflectedNewMethod.Type().NumIn()
	argumentsToCallNewMethod := make([]reflect.Value, constructorInputCount)

	var errors []error
	for i := 0; i < constructorInputCount; i++ {
		reflectedNewMethodArgument := reflectedNewMethod.Type().In(i)

		dependencyName := newMethodArgumentNames[i]

		dependencyFromContainer, err := container.GetSecure(dependencyName, isCached)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		reflectedDependencyFromContainer := reflect.ValueOf(dependencyFromContainer)
		reflectedDependencyFromContainer = replaceCompatibleNilDependency(
			reflectedNewMethodArgument,
			reflectedDependencyFromContainer,
			dependencyFromContainer,
		)

		err = assertCompatible(
			reflectedNewMethodArgument,
			reflect.TypeOf(dependencyFromContainer),
			dependencyName,
			serviceId,
		)
		if err != nil {
			errors = append(errors, err)
		} else {
			argumentsToCallNewMethod[i] = reflectedDependencyFromContainer
		}
	}

	return argumentsToCallNewMethod, mergeErrors(errors)
}

//replaceCompatibleNilDependency will return correct reflect value if dependency from container is nil and it is
//compatible with the interface or pointer constructor argument, otherwise it will return the original reflected
//dependency
func replaceCompatibleNilDependency(
	reflectedConstructorArgument reflect.Type,
	reflectedDependencyFromContainer reflect.Value,
	dependencyFromContainer interface{},
) reflect.Value {
	if dependencyFromContainer == nil && (reflectedConstructorArgument.Kind() == reflect.Interface || reflectedConstructorArgument.Kind() == reflect.Ptr) {
		return reflect.New(reflectedConstructorArgument).Elem()
	}

	return reflectedDependencyFromContainer
}
