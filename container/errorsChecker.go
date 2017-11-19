package container

import "strings"

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicIfErrors(errors []error ) {
	errorStrings := []string{}
	for _, err := range errors {
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}

	if len(errorStrings) > 0  {
		panic(strings.Join(errorStrings, ";\n"))
	}
}
