package tests

import (
	"testing"
	"github.com/breathbath/gotainer/examples"
	"github.com/breathbath/gotainer/container"
)

func TestLazyCollectionDependencies(t *testing.T) {
	cont := examples.CreateContainer()

	cont.AddNoArgumentsConstructor("someMap", func() (interface{}, error){
		return map[string]string{"someMapKey": "someMapValue"}, nil
	})

	cont.AddNoArgumentsConstructor("someMapPointer", func() (interface{}, error){
		return &map[string]string{"someMapKeyPointer": "someMapPointerValue"}, nil
	})

	cont.AddNoArgumentsConstructor("someSlice", func() (interface{}, error){
		return []string{"someSliceValue1", "someSliceValue2"}, nil
	})

	cont.AddNoArgumentsConstructor("someSlicePointer", func() (interface{}, error){
		return &[]string{"someSlicePointerValue3", "someSlicePointerValue4"}, nil
	})

	cont.AddConstructor("someMapFetcher", func(c container.Container) (interface{}, error){
		var mapDependency map[string]string
		c.GetTypedService("someMap", &mapDependency)
		return mapDependency["someMapKey"], nil
	})

	var result string
	cont.GetTypedService("someMapFetcher", &result)

	if result != "someMapValue" {
		t.Errorf("Wrong map fetched from the container for service 'someMapFetcher'", )
	}
}

func NewBooksAuthorMap() map[string] string {
	return map[string] string {"Author1" : "Book1", "Author2": "Book2"}
}

func NewBooksPriceMap() *map[string] int {
	return &map[string] int {"Book1" : 100, "Author2": 150}
}