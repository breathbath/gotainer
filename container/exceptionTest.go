package container

import "testing"

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
		t.Fatalf("Wrong panic message: '%s', expected message: '%s'", panicMessage, expectedPanicName)
	}
}
