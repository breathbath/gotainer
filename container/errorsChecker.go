package container

import (
	coreErrors "errors"
	"strings"
	"testing"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func mergeErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	errorStrings := []string{}
	for _, err := range errors {
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}

	return coreErrors.New(strings.Join(errorStrings, ";\n"))
}

func assertErrorText(expectedErrorText string, providedError error, t *testing.T) {
	if providedError == nil {
		t.Errorf("Error '%s' was expected but none was returned", expectedErrorText)
		return
	}

	if providedError.Error() != expectedErrorText {
		unexpectedPanicMessage := "\nWrong error text:(-expected text, +provided text)\n- %s\n+ %s"
		t.Errorf(unexpectedPanicMessage, expectedErrorText, providedError.Error())
	}
}

func assertNoError(providedError error, t *testing.T) {
	if providedError == nil {
		return
	}

	t.Errorf("Unexpected error: %v", providedError)
}
