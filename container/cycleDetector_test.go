package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

var cycleTree Tree

func init() {
	cycleTree = Tree{
		Node{
			ID:           "userProvider",
			NewFunc:      mocks.NewUserProvider,
			ServiceNames: Services{"roleProvider"},
		},
		Node{
			ID:           "roleProvider",
			NewFunc:      mocks.NewRoleProvider,
			ServiceNames: Services{"userProvider"},
		},
		Node{
			ID:      "rightsProvider",
			NewFunc: mocks.NewRightsProvider,
		},
	}
}

func TestSelfReferenceFailureWithConfigDeclaration(t *testing.T) {
	defer ExpectPanic(t, "Detected dependencies' cycle: book_storage->book_storage [check 'book_storage' service]")
	recursiveTree := Tree{
		Node{
			ID:           "book_storage",
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
	defer ExpectPanic(
		t,
		"Detected dependencies' cycle: userProvider->roleProvider->userProvider [check 'roleProvider' service] [check 'userProvider' service]",
	)

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
		"Detected dependencies' cycle: rolesProvider->userProvider->rolesProvider [check 'userProvider' service] [check 'rolesProvider' service]",
		"Detected dependencies' cycle: userProvider->rolesProvider->userProvider [check 'userProvider' service] [check 'rolesProvider' service]",
	)
	cont := NewRuntimeContainer()

	cont.AddConstructor("rolesProvider", func(c Container) (interface{}, error) {
		return c.GetSecure("userProvider", true)
	})

	cont.AddConstructor("userProvider", func(c Container) (interface{}, error) {
		return c.GetSecure("rolesProvider", true)
	})

	cont.Get("rolesProvider", true)
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

func TestCycledAndNonCycledDependenciesWithSecureMethod(t *testing.T) {
	cont := NewRuntimeContainer()
	cont.AddConstructor("configPath", func(c Container) (interface{}, error) {
		return c.GetSecure("configPath", true)
	})
	cont.AddConstructor("email", func(c Container) (interface{}, error) {
		return "root@root.me", nil
	})

	cont.GetSecure("configPath", true) //is cyclic we just ignore it

	email, err := cont.GetSecure("email", true)
	if err != nil {
		t.Errorf("No error is expected for a non-cyclic dependency but %v received", err)
		return
	}
	if email.(string) != "root@root.me" {
		t.Errorf("Should get 'root@root.me' value for email dependency but got %s", email)
		return
	}
}

func TestResetCycleDetector(t *testing.T) {
	cd := NewCycleDetector()
	cd.cycleDetected = true
	if !cd.HasCycle() {
		t.Error("Cycle flag should set as detected")
		return
	}

	cd.Reset()

	if cd.HasCycle() {
		t.Error("Cycle flag should set as not detected after reset")
		return
	}

	cd.cycleDetected = true
	cd.DisableCycleDetection()
	cd.Reset()
	if !cd.HasCycle() {
		t.Error("Cycle flag should set as detected if cycle detection is disabled")
		return
	}
}

func TestCycleDetectionWithCache(t *testing.T) {
	cont := NewRuntimeContainer()

	cont.AddConstructor("pathRegistry", func(c Container) (interface{}, error) {
		return "pathRegistry", nil
	})
	cont.AddDependencyObserver("newPath", "pathRegistry", func(pathRegistry, currentPath string) {
	})

	cont.AddConstructor("rootPath", func(c Container) (interface{}, error) {
		return "root", nil
	})
	cont.AddConstructor("imgPath", func(c Container) (interface{}, error) {
		rootPath := c.Get("rootPath", true).(string)
		return rootPath + "/" + "img", nil
	})
	cont.RegisterDependencyEvent("newPath", "imgPath")

	cont.AddConstructor("configPath", func(c Container) (interface{}, error) {
		rootPath := c.Get("rootPath", true).(string)
		return rootPath + "/" + "config", nil
	})
	cont.RegisterDependencyEvent("newPath", "configPath")

	cont.AddConstructor("envPath", func(c Container) (interface{}, error) {
		configPath := c.Get("configPath", true).(string)
		return configPath + "/" + "env", nil
	})
	cont.RegisterDependencyEvent("newPath", "envPath")

	cont.Get("configPath", true)

	cont.Get("pathRegistry", true)
}

func TestCycleDetectionWithEvents(t *testing.T) {
	cont := NewRuntimeContainer()
	cont.AddConstructor("pathsCollector", func(c Container) (interface{}, error) {
		return &mocks.PathsCollector{Paths: []string{}}, nil
	})

	cont.AddDependencyObserver("newPathProvider", "pathsCollector", func(pc *mocks.PathsCollector, newPathProvider mocks.PathProvider) {
		pc.AddPath(newPathProvider.GetPath())
	})

	cont.AddConstructor("longestPathProvider", func(c Container) (interface{}, error) {
		return &mocks.LongestPathProvider{}, nil
	})

	cont.AddDependencyObserver("possibleLongestPathProvider", "longestPathProvider", func(lpp *mocks.LongestPathProvider, newPathProvider mocks.SimplePathProvider) {
		lpp.EvaluatePath(newPathProvider.GetPath())
	})

	cont.AddConstructor("pathAProvider", func(c Container) (interface{}, error) {
		return mocks.SimplePathProvider{Path: "pathA"}, nil
	})
	cont.RegisterDependencyEvent("newPathProvider", "pathAProvider")
	cont.RegisterDependencyEvent("possibleLongestPathProvider", "pathAProvider")

	cont.AddConstructor("pathBProvider", func(c Container) (interface{}, error) {
		return mocks.SimplePathProvider{Path: "pathB"}, nil
	})
	cont.RegisterDependencyEvent("newPathProvider", "pathBProvider")
	cont.RegisterDependencyEvent("possibleLongestPathProvider", "pathBProvider")

	cont.AddConstructor("pathCProvider", func(c Container) (interface{}, error) {
		return mocks.SimplePathProvider{Path: "pathCC"}, nil
	})
	cont.RegisterDependencyEvent("newPathProvider", "pathCProvider")
	cont.RegisterDependencyEvent("possibleLongestPathProvider", "pathCProvider")

	cont.RegisterDependencyEvent("newPathProvider", "longestPathProvider")

	pc := cont.Get("pathsCollector", true).(*mocks.PathsCollector)
	providedPaths := pc.GetAllPaths()
	expectedPaths := "pathA,pathB,pathCC,pathCC"
	if providedPaths != expectedPaths {
		t.Errorf("Not expected result %s is returned by pathsCollector, expected result was: %s", providedPaths, expectedPaths)
	}
}
