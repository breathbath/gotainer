package container

import (
	coreErrors "errors"
	"strings"
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
