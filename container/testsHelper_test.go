package container

import (
	"testing"
	"fmt"
)

//ExpectPanic a helper method to simulate a  panic expectation in tests
func ExpectPanic(expectedPanicName string, t *testing.T) {
	err := recover()
	if err == nil {
		t.Fatalf("There was a panic message expected '%s', but none was received", expectedPanicName)
	}
	var panicMessage string

	switch errorName := err.(type) {
	case string:
		panicMessage = errorName
	case error:
		panicMessage = errorName.Error()
	default:
		panicMessage = "Unknown error type"
	}

	if panicMessage != expectedPanicName {
		t.Fatalf("\nWrong panic message:(-expected message, +provided message)\n- %s\n+ %s", expectedPanicName, panicMessage)
	}
}

func AssertExpectedDependency(c Container, dependencyName string, expectedValue interface{}, t *testing.T) {
	result := c.Get(dependencyName, true)
	if result != expectedValue {
		t.Errorf(
			"Unexpected result for service '%s': received result is '%s', expected value is '%s'",
			dependencyName,
			fmt.Sprint(result),
			fmt.Sprint(expectedValue),
		)
	}
}

func AssertError(err error, expectedError string,  t *testing.T) {
	if err == nil {
		t.Errorf("Error '%s' is expected but non was returned", expectedError)
		return
	}

	if err.Error() != expectedError {
		t.Errorf("Error '%s' is expected but an '%s' error was returned", expectedError, err.Error())
	}
}
