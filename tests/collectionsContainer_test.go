package tests

import (
	"errors"
	"github.com/breathbath/gotainer/container"
	"github.com/breathbath/gotainer/examples"
	"testing"
)

func TestLazyCollectionDependencies(t *testing.T) {
	cont := examples.CreateContainer()

	cont.AddConstructor("someMap", func(c container.Container) (interface{}, error) {
		return map[string]string{"someMapKey": "someMapValue"}, nil
	})

	cont.AddConstructor("someSlice", func(c container.Container) (interface{}, error) {
		return []string{"someSliceValue1", "someSliceValue2"}, nil
	})

	AssertStringValueExtracted(
		"someMapValue",
		"someMapFetcher",
		func(c container.Container) (interface{}, error) {
			var mapDependency map[string]string
			c.Scan("someMap", &mapDependency)
			return mapDependency["someMapKey"], nil
		},
		cont,
		t,
	)

	AssertStringValueExtracted(
		"someSliceValue2",
		"someSliceFetcher",
		func(c container.Container) (interface{}, error) {
			var sliceString []string
			c.Scan("someSlice", &sliceString)
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
	cont.AddNewMethod("book_prices", examples.GetBookPrices)
	cont.AddNewMethod("books", examples.GetAllBooks)
	cont.AddNewMethod("price_finder", examples.NewBooksPriceFinder, "book_prices", "books")

	var priceCalculator func(bookId string) int
	cont.Scan("price_finder", &priceCalculator)

	expectedResult := 100
	result := priceCalculator("1")
	if priceCalculator("1") != expectedResult {
		t.Errorf("Wrong price calculator result '%d', expected result is '%d'", result, expectedResult)
	}
}

func AssertStringValueExtracted(expectedString string, extractFuncName string, extractFunc container.Constructor, c container.Container, t *testing.T) {
	c.AddConstructor(extractFuncName, extractFunc)
	var result string
	c.Scan(extractFuncName, &result)
	if result != expectedString {
		t.Errorf("Unexpected string '%s' fetched from the container for dependency '%s'. Expected string was '%s'", result, extractFuncName, expectedString)
	}
}
