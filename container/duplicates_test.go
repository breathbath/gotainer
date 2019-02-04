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

	defer ExpectPanic(
		t,
		"Detected duplicated dependency declaration 'bookStorage'",
	)

	RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfigWithDuplicates)
}

func TestConfigConstrDuplicate(t *testing.T) {
	exampleConfigWithDuplicates := append(exampleConfig, Node{
		ID: "db",
		Constr: func(c Container) (interface{}, error) {
			return "", nil
		},
	})

	defer ExpectPanic(
		t,
		"Detected duplicated dependency declaration 'db'",
	)

	RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfigWithDuplicates)
}

func TestParamConfigDuplicates(t *testing.T) {
	exampleConfigWithDuplicates := append(exampleConfig, Node{
		Parameters: map[string]interface{}{
			"connectionString": "someOtherStr",
		},
	})

	defer ExpectPanic(
		t,
		"Detected duplicated dependency declaration 'connectionString'",
	)

	RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfigWithDuplicates)
}

func TestNewMethodDuplicates(t *testing.T) {
	defer ExpectPanic(
		t,
		"Detected duplicated dependency declaration 'bookShelve'",
	)

	cont := NewRuntimeContainer()

	cont.AddNewMethod("bookShelve", mocks.NewBookShelve)
	cont.AddNewMethod("bookShelve", mocks.NewConfig)
}

func TestConstrAndNewMethodDuplicates(t *testing.T) {
	defer ExpectPanic(
		t,
		"Detected duplicated dependency declaration 'bookShelve'",
	)

	cont := NewRuntimeContainer()

	cont.AddNewMethod("bookShelve", mocks.NewBookShelve)
	cont.AddConstructor("bookShelve", func(c Container) (interface{}, error) {
		return "", nil
	})
}

func TestHybridConfigWithDirectDuplicatesDeclaration(t *testing.T) {
	defer ExpectPanic(
		t,
		"Detected duplicated dependency declaration 'db'",
	)

	cont, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfig)
	if err != nil {
		t.Error(err)
		return
	}

	cont.AddNewMethod("db", mocks.NewBookShelve)
}

func TestMergeDuplicates(t *testing.T) {
	defer ExpectPanic(
		t,
		"Cannot merge containers because of non unique Service id 'bookShelve'",
	)

	cont1 := NewRuntimeContainer()
	cont1.AddNewMethod("bookShelve", mocks.NewBookShelve)

	cont2 := NewRuntimeContainer()
	cont2.AddNewMethod("bookShelve", mocks.NewConfig)

	cont1.Merge(cont2)
}
