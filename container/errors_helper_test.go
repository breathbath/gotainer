package container

import "testing"

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
