package container

import (
	"fmt"
	"strings"
	"testing"
)

//ExpectPanic a helper method to simulate a panic expectation in tests, you can provide multiple possible panic variants
func ExpectPanic( t *testing.T, expectedPanicNames ...string) {
	expectedPanicNamesFlat := strings.Join(expectedPanicNames, ",")

	noPanicDetectedMessage := "There was a panic message expected '%s', but none was received"
	unexpectedPanicMessage := "\nWrong panic message:(-expected message, +provided message)\n- %s\n+ %s"
	if len(expectedPanicNames) > 1 {
		noPanicDetectedMessage = "There were panic variants expected [%s], but none was received"
		unexpectedPanicMessage = "\nWrong panic message:(-expected message variants, +provided message)\n- %s\n+ %s"
	}

	err := recover()
	if err == nil {
		if len(expectedPanicNames) == 0 {
			return
		}

		t.Fatalf(noPanicDetectedMessage, expectedPanicNamesFlat)
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

	for _, expectedPanicName := range expectedPanicNames {
		if panicMessage == expectedPanicName {
			return
		}
	}

	t.Fatalf(unexpectedPanicMessage, expectedPanicNamesFlat, panicMessage)
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

func AssertError(err error, expectedError string, t *testing.T) {
	if err == nil {
		t.Errorf("Error '%s' is expected but non was returned", expectedError)
		return
	}

	if err.Error() != expectedError {
		t.Errorf("Error '%s' is expected but an '%s' error was returned", expectedError, err.Error())
	}
}
