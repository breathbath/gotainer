package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

var exampleConfig Tree

func init() {
	exampleConfig = Tree{
		Node{
			Parameters: map[string]interface{}{
				"connectionString": "someStr",
			},
		},
		Node{
			ID:           "db",
			NewFunc:      mocks.NewFakeDb,
			ServiceNames: Services{"connectionString"},
		},
		Node{
			ID:           "bookFinder",
			NewFunc:      mocks.NewBookFinder,
			ServiceNames: Services{"bookStorage", "bookCreator",},
		},
		Node{
			ID:           "bookStorage",
			NewFunc:      mocks.NewBookStorage,
			ServiceNames: Services{"db"},
		},
		Node{
			ID: "bookCreator",
			Constr: func(c Container) (interface{}, error) {
				return mocks.BookCreator{}, nil
			},
		},
	}
}

func TestConfigNewMethodDuplicate(t *testing.T) {
	exampleConfigWithDuplicates := append(exampleConfig, Node{
		ID:      "bookStorage",
		NewFunc: mocks.NewInMemoryCache,
	})

	_, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfigWithDuplicates)
	AssertError(err, "Detected duplicated dependency declaration 'bookStorage'", t)
}

func TestConfigConstrDuplicate(t *testing.T) {
	exampleConfigWithDuplicates := append(exampleConfig, Node{
		ID: "db",
		Constr: func(c Container) (interface{}, error) {
			return "", nil
		},
	})

	_, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfigWithDuplicates)
	AssertError(err, "Detected duplicated dependency declaration 'db'", t)
}

func TestParamConfigDuplicates(t *testing.T) {
	exampleConfigWithDuplicates := append(exampleConfig, Node{
		Parameters: map[string]interface{}{
			"connectionString": "someOtherStr",
		},
	})

	_, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfigWithDuplicates)
	if err != nil {
		AssertError(err, "Detected duplicated dependency declaration 'connectionString'", t)
	}
}

func TestNewMethodDuplicates(t *testing.T) {
	cont := NewRuntimeContainer()

	err := cont.AddNewMethod("bookShelve", mocks.NewBookShelve)
	assertNoError(err, t)

	err = cont.AddNewMethod("bookShelve", mocks.NewConfig)
	AssertError(err, "Detected duplicated dependency declaration 'bookShelve'", t)
}

func TestConstrAndNewMethodDuplicates(t *testing.T) {
	cont := NewRuntimeContainer()

	err := cont.AddNewMethod("bookShelve", mocks.NewBookShelve)
	assertNoError(err, t)

	err = cont.AddConstructor("bookShelve", func(c Container) (interface{}, error) {
		return "", nil
	})
	AssertError(err, "Detected duplicated dependency declaration 'bookShelve'", t)
}

func TestHybridConfigWithDirectDuplicatesDeclaration(t *testing.T) {
	cont, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfig)
	if err != nil {
		t.Error(err)
		return
	}

	err = cont.AddNewMethod("db", mocks.NewBookShelve)
	AssertError(err, "Detected duplicated dependency declaration 'db'", t)
}

func TestMergeDuplicates(t *testing.T) {
	cont1 := NewRuntimeContainer()
	cont1.AddNewMethod("bookShelve", mocks.NewBookShelve)

	cont2 := NewRuntimeContainer()
	cont2.AddNewMethod("bookShelve", mocks.NewConfig)

	err := cont1.Merge(cont2)
	AssertError(err, "Cannot merge containers because of non unique Service id 'bookShelve'", t)
}
