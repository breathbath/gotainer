package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestParametersAdding(t *testing.T) {
	c := CreateContainer()

	stringParams := map[string]string{
		"paramString1": "valueString1",
		"paramString2": "valueString2",
	}

	boolParams := map[string]bool{"paramBool1": true}
	intParams := map[string]int{"paramInt1": 2}
	int64Params := map[string]int64{"paramInt641": 3}

	bookParams := map[string]mocks.Book{
		"paramBook1": {Id: "book1"},
		"paramBook2": {Id: "book2"},
	}

	stringPointer := "valueStringPointer1"
	pointerParams := map[string]*string{"paramStringPointer1": &stringPointer}

	RegisterParameters(c, stringParams, boolParams, intParams, int64Params, bookParams, pointerParams)

	AssertExpectedDependency(c, "paramString1", "valueString1", t)
	AssertExpectedDependency(c, "paramString2", "valueString2", t)
	AssertExpectedDependency(c, "paramBool1", true, t)
	AssertExpectedDependency(c, "paramInt1", 2, t)
	AssertExpectedDependency(c, "paramInt641", int64(3), t)

	var book1, book2 mocks.Book
	c.Scan("paramBook1", &book1)
	c.Scan("paramBook2", &book2)

	assertExpectedBook("paramBook1", "book1", book1, t)
	assertExpectedBook("paramBook2", "book2", book2, t)

	AssertExpectedDependency(c, "paramStringPointer1", &stringPointer, t)
}

func TestIgnoringNilValues(t *testing.T) {
	c := CreateContainer()
	nilParams := map[string]map[int]int{
		"nilParam": {},
	}
	RegisterParameters(c, nilParams)
}

func TestFailingForNonMapInput(t *testing.T) {
	c := CreateContainer()
	err := RegisterParameters(c, "some_string")
	AssertError(err, "A map type should be provided to register parameters", t)
}

func TestFailingForNonStringMapKeys(t *testing.T) {
	c := CreateContainer()
	err := RegisterParameters(c, map[int]int{1: 22})
	AssertError(err, "A map[string]interface{} should be provided to register parameters", t)
}

func assertExpectedBook(serviceId, expectedBookId string, book mocks.Book, t *testing.T) {
	if book.Id != expectedBookId {
		t.Errorf(
			"Unexpected book value for %s, received book with id %s, expected book %s",
			serviceId,
			book.Id,
			expectedBookId,
		)
	}
}
