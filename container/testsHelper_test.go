package container

import (
	"fmt"
	"strings"
	"testing"
)

//ExpectPanic a helper method to simulate a panic expectation in tests, you can provide multiple possible panic variants
func ExpectPanic(t *testing.T, expectedPanicNames ...string) {
	ExpectErrorVariant(recover(), t, expectedPanicNames...)
}

func ExpectErrorSubmatch(err error, expectedSubmatch string, t *testing.T, ) {
	if strings.Contains(err.Error(), expectedSubmatch) {
		return
	}
	t.Errorf("The returned error '%s' is expected to contain '%s' but it does not", err.Error(), expectedSubmatch)
}

//ExpectErrorVariant a helper method to check errors in tests, you can provide multiple possible error variants
func ExpectErrorVariant(err interface{}, t *testing.T, expectedErrorMessages ...string) {
	expectedPanicNamesFlat := strings.Join(expectedErrorMessages, ",")

	noPanicDetectedMessage := "There was a panic message expected '%s', but none was received"
	unexpectedPanicMessage := "\nWrong panic message:(-expected message, +provided message)\n- %s\n+ %s"
	if len(expectedErrorMessages) > 1 {
		noPanicDetectedMessage = "There were panic variants expected [%s], but none was received"
		unexpectedPanicMessage = "\nWrong panic message:(-expected message variants, +provided message)\n- %s\n+ %s"
	}

	if err == nil {
		if len(expectedErrorMessages) == 0 {
			return
		}

		t.Errorf(noPanicDetectedMessage, expectedPanicNamesFlat)
	}

	var errorMessage string

	switch errorName := err.(type) {
	case string:
		errorMessage = errorName
	case error:
		errorMessage = errorName.Error()
	default:
		errorMessage = "Unknown error type"
	}

	for _, expectedPanicName := range expectedErrorMessages {
		if errorMessage == expectedPanicName {
			return
		}
	}

	t.Errorf(unexpectedPanicMessage, expectedPanicNamesFlat, errorMessage)
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
