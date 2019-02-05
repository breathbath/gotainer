package container

import (
	"errors"
	"testing"
)

func TestMergeErrors(t *testing.T) {
	errs := []error{}

	mergedErr := mergeErrors(errs)
	if mergedErr != nil {
		t.Errorf("Nil error is expected after merge of empty errors slice")
	}

	errs = append(errs, errors.New("Some error"), errors.New("Some other error"))
	mergedErr = mergeErrors(errs)
	if mergedErr == nil {
		t.Errorf("It is an error expected but non is returned")
	}

	AssertError(mergedErr, "Some error;\nSome other error", t)
}
