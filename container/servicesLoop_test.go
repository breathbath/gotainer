package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestSelfReferenceFailureWithConfigDeclaration(t *testing.T) {
	defer ExpectPanic("Recursive self reference declaration [check 'book_storage' service]", t)
	recursiveTree := Tree{
		Node{
			Id:           "book_storage",
			NewFunc:      mocks.NewBookStorage,
			ServiceNames: Services{"book_storage"},
		},
	}

	cont := RuntimeContainerBuilder{}.BuildContainerFromConfig(recursiveTree)
	cont.Check()
}

func TestSelfReferenceFailuresWitDirectDeclaration(t *testing.T) {
	defer ExpectPanic("Recursive self reference declaration [check 'book_finder' service]", t)
	cont := NewRuntimeContainer()

	cont.AddConstructor("db", func(c Container) (interface{}, error) {
		return mocks.NewFakeDb("someConnectionString"), nil
	})
	cont.AddNewMethod("book_storage", mocks.NewBookStorage, "db")
	cont.AddNewMethod("book_finder", mocks.NewBookFinder, "book_storage", "book_finder")

	cont.Check()
}

func TestCircleReferencesWithConfigDeclaration(t *testing.T) {
	circleTree := Tree{
		Node{
			Id:           "userProvider",
			NewFunc:      mocks.NewUserProvider,
			ServiceNames: Services{"roleProvider"},
		},
		Node{
			Id:           "roleProvider",
			NewFunc:      mocks.NewRoleProvider,
			ServiceNames: Services{"userProvider"},
		},
		Node{
			Id:           "rightsProvider",
			NewFunc:      mocks.NewRightsProvider,
		},
	}

	cont := RuntimeContainerBuilder{}.BuildContainerFromConfig(circleTree)
	cont.Check()
}

func TestCircleReferencesWithDirectDeclaration(t *testing.T) {
	cont := NewRuntimeContainer()

	cont.AddNewMethod("userProvider", mocks.NewUserProvider, "roleProvider")
	cont.AddNewMethod("roleProvider", mocks.NewRoleProvider, "userProvider")

	cont.Check()
}

func TestCircleReferencesWithConstructor(t *testing.T) {
	cont := NewRuntimeContainer()

	cont.AddConstructor("userProvider", func(c Container) (interface{}, error) {
		rp := c.Get("roleProvider", true).(mocks.RoleProvider)
		return rp, nil
	})

	cont.AddConstructor("roleProvider", func(c Container) (interface{}, error) {
		rp := c.Get("userProvider", true).(mocks.UserProvider)
		return rp, nil
	})

	cont.Check()
}