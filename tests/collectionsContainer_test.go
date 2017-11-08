package tests

import (
	"testing"
	"github.com/breathbath/gotainer/examples"
	"github.com/breathbath/gotainer/container"
	"errors"
)

func TestLazyCollectionDependencies(t *testing.T) {
	cont := examples.CreateContainer()

	cont.AddNoArgumentsConstructor("someMap", func() (interface{}, error){
		return map[string]string{"someMapKey": "someMapValue"}, nil
	})

	cont.AddNoArgumentsConstructor("someSlice", func() (interface{}, error){
		return []string{"someSliceValue1", "someSliceValue2"}, nil
	})

	AssertStringValueExtracted(
		"someMapValue",
		"someMapFetcher",
		func(c container.Container) (interface{}, error){
			var mapDependency map[string]string
			c.GetTypedService("someMap", &mapDependency)
			return mapDependency["someMapKey"], nil
		},
		cont,
		t,
	)

	AssertStringValueExtracted(
		"someSliceValue2",
		"someSliceFetcher",
		func(c container.Container) (interface{}, error){
			var sliceString []string
			c.GetTypedService("someSlice", &sliceString)
			if len(sliceString) < 2 {
				return "", errors.New("Slice with 2 values is expected for 'someSlice' dependency, but none is returned in 'someSliceFetcher'")
			}
			return sliceString[1], nil
		},
		cont,
		t,
	)
}

func TestCollectionDependencies(t *testing.T) {
	cont := examples.CreateContainer()
	cont.AddTypedConstructor("book_prices", examples.GetBookPrices)
	cont.AddTypedConstructor("books", examples.GetAllBooks)
	cont.AddTypedConstructor("price_finder", examples.NewBooksPriceFinder, "book_prices", "books")

	var priceCalculator func(bookId string) int
	cont.GetTypedService("price_finder", &priceCalculator)

	expectedResult := 100
	result := priceCalculator("1")
	if priceCalculator("1") != expectedResult {
		t.Errorf("Wrong price calculator result '%d', expected result is '%d'", result, expectedResult)
	}
}

func AssertStringValueExtracted(expectedString string, extractFuncName string, extractFunc container.ArgumentsConstructor, c container.Container, t *testing.T) {
	c.AddConstructor(extractFuncName, extractFunc)
	var result string
	c.GetTypedService(extractFuncName, &result)
	if result != expectedString {
		t.Errorf("Unexpected string '%s' fetched from the container for service '%s'. Expected string was '%s'", result, extractFuncName, expectedString)
	}
}