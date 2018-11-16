package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

var cycleTree Tree

func init() {
	cycleTree = Tree{
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
			Id:      "rightsProvider",
			NewFunc: mocks.NewRightsProvider,
		},
	}
}

func TestSelfReferenceFailureWithConfigDeclaration(t *testing.T) {
	defer ExpectPanic(t, "Detected dependencies' cycle: book_storage->book_storage [check 'book_storage' service]")
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
	defer ExpectPanic(t, "Detected dependencies' cycle: book_finder->book_finder [check 'book_finder' service]")
	cont := NewRuntimeContainer()

	cont.AddConstructor("db", func(c Container) (interface{}, error) {
		return mocks.NewFakeDb("someConnectionString"), nil
	})
	cont.AddNewMethod("book_storage", mocks.NewBookStorage, "db")
	cont.AddNewMethod("book_finder", mocks.NewBookFinder, "book_storage", "book_finder")

	cont.Get("book_finder", true)
}

func TestCycleReferencesWithConfigDeclaration(t *testing.T) {
	defer ExpectPanic(t, "Detected dependencies' cycle: userProvider->roleProvider->userProvider [check 'roleProvider' service] [check 'userProvider' service]")

	cont := RuntimeContainerBuilder{}.BuildContainerFromConfig(cycleTree)
	cont.Get("userProvider", true)
}

func TestCycleReferencesWithNewMethodDeclaration(t *testing.T) {
	defer ExpectPanic(
		t,
		"Detected dependencies' cycle: userProvider->roleProvider->userProvider [check 'roleProvider' service] [check 'userProvider' service]",
		"Detected dependencies' cycle: roleProvider->userProvider->roleProvider [check 'userProvider' service] [check 'roleProvider' service]",
	)

	cont := NewRuntimeContainer()

	cont.AddNewMethod("userProvider", mocks.NewUserProvider, "roleProvider")
	cont.AddNewMethod("roleProvider", mocks.NewRoleProvider, "userProvider")

	cont.Check()
}

func TestCycleReferencesWithConstructor(t *testing.T) {
	defer ExpectPanic(
		t,
		"Detected dependencies' cycle: userReader->nameCutter->userReader",
		"Detected dependencies' cycle: nameCutter->userReader->nameCutter",
	)
	cont := NewRuntimeContainer()

	cont.AddConstructor("nameCutter", func(c Container) (interface{}, error) {
		rp := c.Get("userReader", true).(mocks.RoleProvider)
		return rp, nil
	})

	cont.AddConstructor("userReader", func(c Container) (interface{}, error) {
		rp := c.Get("nameCutter", true).(mocks.UserProvider)
		return rp, nil
	})

	cont.Check()
}

func TestNoCycleWithMultipleReferencedDependencies(t *testing.T) {
	cont := NewRuntimeContainer()
	cont.AddConstructor("dbConnector", func(c Container) (interface{}, error) {
		usrPass := c.Get("dbUser", true).(string)
		dbPass := c.Get("dbPass", true).(string)
		return usrPass + dbPass, nil
	})

	cont.AddConstructor("dbUser", func(c Container) (interface{}, error) {
		c.Get("configPath", true)
		return "root", nil
	})

	cont.AddConstructor("configPath", func(c Container) (interface{}, error) {
		return "/temp", nil
	})

	cont.AddConstructor("dbPass", func(c Container) (interface{}, error) {
		c.Get("configPath", true)
		return "rootpass", nil
	})

	cont.Get("dbConnector", true)
	cont.Check()
}

func TestNoCycleWithMultipleCallsOfSameDependency(t *testing.T) {
	cont := NewRuntimeContainer()
	cont.AddConstructor("configPath", func(c Container) (interface{}, error) {
		return "/tmp", nil
	})

	cont.Get("configPath", true)
	cont.Get("configPath", true)
	cont.Check()
}

func TestCyclicAndNonCyclicDependencies(t *testing.T) {
	defer ExpectPanic(
		t,
		"Detected dependencies' cycle: configPath->configPath",
	)

	cont := NewRuntimeContainer()
	cont.AddConstructor("configPath", func(c Container) (interface{}, error) {
		return c.Get("configPath", true), nil
	})
	cont.AddConstructor("email", func(c Container) (interface{}, error) {
		return "root@root.me", nil
	})

	email := cont.Get("email", true).(string)
	if email != "root@root.me" {
		t.Errorf("Should get 'root@root.me' value for email dependency but got %s", email)
		return
	}
	cont.Get("configPath", true)
}

func TestCycleDetectionWithSecureMethod(t *testing.T) {
	cont := RuntimeContainerBuilder{}.BuildContainerFromConfig(cycleTree)
	_, err := cont.GetSecure("userProvider", true)
	expectedErrorText := "Detected dependencies' cycle: userProvider->roleProvider->userProvider [check 'roleProvider' service] [check 'userProvider' service]"
	if err.Error() != expectedErrorText {
		t.Errorf("Error %s expected but %s was received", expectedErrorText, err.Error())
	}
}