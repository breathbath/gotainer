package container

import (
	"errors"
	"fmt"
)

// This example fetching of services from container with a wrapped error
// if it happens
func ExampleFetch_ScanSecure() {
	cont := NewRuntimeContainer()
	cont.AddConstructor("parameterOk", func(c Container) (interface{}, error) {
		return "123456", nil
	})
	cont.AddConstructor("parameterFailed", func(c Container) (interface{}, error) {
		return "123456", errors.New("Some error")
	})

	var parameterOk string
	err := cont.ScanSecure("parameterOk", true, &parameterOk)
	fmt.Println(err)

	err = cont.ScanSecure("parameterFailed", true, &parameterOk)
	fmt.Println(err)

	err = cont.ScanSecure("unknownParameter", true, &parameterOk)
	fmt.Println(err)

	// Output:
	// <nil>
	// Some error [check 'parameterFailed' service]
	// Unknown dependency 'unknownParameter'
}