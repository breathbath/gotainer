package container

import (
	"reflect"
	"errors"
	"fmt"
)

func copySurceVariableToDestinationVariable(sourceVariable interface{}, destination interface{}, serviceName string) error {
	dpv := reflect.ValueOf(destination)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value")
	}
	if dpv.IsNil() {
		return errors.New("nil pointer passed to destination")
	}

	sv := reflect.ValueOf(sourceVariable)

	dv := reflect.Indirect(dpv)
	dvKind := dv.Kind()
	svKind := sv.Kind()
	if dvKind == svKind && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	errStr := fmt.Sprintf(
		"Cannot convert created value of type '%s' to expected destination value '%s' for sourceVariable declaration %s",
		sv.Type().Name(),
		dv.Type().Name(),
		serviceName,
	)

	return errors.New(errStr)
}

func isErrorType(reflectedType reflect.Type) bool {
	expectedErrorInterface := reflect.TypeOf((*error)(nil)).Elem()
	return reflectedType.Implements(expectedErrorInterface)
}
