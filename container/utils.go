package container

import (
	"errors"
	"fmt"
	"reflect"
	"github.com/breathbath/gotainer/container/mocks"
)

//copySourceVariableToDestinationVariable copies a dependency fetched from the container
//to the pointer reference provided as dest in Scan or ScanNonCached of the container
func copySourceVariableToDestinationVariable(createdDependency interface{}, destination interface{}, dependencyName string) error {
	reflectedDestinationValue := reflect.ValueOf(destination)
	if reflectedDestinationValue.Kind() != reflect.Ptr {
		return errors.New("Cannot set value to a non pointer variable")
	}

	reflectedCreatedDependency := reflect.ValueOf(createdDependency)

	if reflectedDestinationValue.IsNil() {
		d := destination.(*mocks.BookShelve)
		va := reflect.ValueOf(&d).Elem()
		v := reflect.New(va.Type().Elem())
		va.Set(v)
		reflectedDestinationValue = reflect.ValueOf(d)
		//va := reflect.ValueOf(&destination).Elem()
		//v := reflect.New(va.Type().Elem())
		//va.Set(v)
		//reflectedDestinationValue = va
		//intType := reflect.TypeOf(&destination).Elem()
		////intPtr2 := reflect.New(intType)
		//
		////indirectedCreatedDependency := reflect.Indirect(reflectedCreatedDependency)
		//indirectedCreatedDependency := reflect.Indirect(reflectedCreatedDependency)
		////indirectedReflectedDestinationValue := intPtr2.Elem()
		//convertedDependencyToDestinationValue := indirectedCreatedDependency.Convert(intType)
		//reflect.ValueOf(destination).Set(convertedDependencyToDestinationValue)
		//return nil
		//rrr := reflect.TypeOf(destination).Elem()
		//rrr2 := reflect.New(rrr)
		//rrr3 := rrr2.Elem().Interface().(mocks.BookShelve)
		//
		//destType := reflect.TypeOf(&destination).Elem()
		//indirectedCreatedDependency := reflect.Indirect(reflectedCreatedDependency)
		//convertedDependencyToDestinationValue := indirectedCreatedDependency.Convert(destType)
		//indirectedReflectedDestinationValue := reflect.Indirect(reflect.ValueOf(destination))
		//indirectedReflectedDestinationValue.Set(convertedDependencyToDestinationValue)

		//reflectedDestinationValue = reflect.ValueOf(destination).Elem()
		//reflectedDestinationValue = reflect.New(reflectedDestinationValue.Type().Elem())
		//indirectedReflectedDestinationValue1 := reflect.Indirect(reflectedDestinationValue)
		//destinationValueType1 := indirectedReflectedDestinationValue1.Type()
		//indirectedReflectedDestinationValue1.Set(reflectedCreatedDependency.Convert(destinationValueType1))
		//return nil
		//fmt.Println(destination)
		//return errors.New("nil pointer passed to destination")
	}

	indirectedReflectedDestinationValue := reflect.Indirect(reflectedDestinationValue)
	destinationValueType := indirectedReflectedDestinationValue.Type()

	if reflectedCreatedDependency.Kind() == reflect.Ptr && sourceCanBeCopiedToDestination(reflectedCreatedDependency, reflectedDestinationValue) {
		indirectedCreatedDependency := reflect.Indirect(reflectedCreatedDependency)
		convertedDependencyToDestinationValue := indirectedCreatedDependency.Convert(destinationValueType)
		indirectedReflectedDestinationValue.Set(convertedDependencyToDestinationValue)
		return nil
	}

	if sourceCanBeCopiedToDestination(reflectedCreatedDependency, indirectedReflectedDestinationValue) {
		indirectedReflectedDestinationValue.Set(reflectedCreatedDependency.Convert(destinationValueType))
		return nil
	}

	errStr := fmt.Sprintf(
		"Cannot convert created value of type '%s' to expected destination value '%s' for createdDependency declaration %s",
		reflectedCreatedDependency.Type().Name(),
		indirectedReflectedDestinationValue.Type().Name(),
		dependencyName,
	)

	return errors.New(errStr)
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
