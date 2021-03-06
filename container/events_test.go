package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestContainerEvents(t *testing.T) {
	c := CreateContainer()
	c.AddDependencyObserver("non_registered_event", "statistics_gateway", func(a interface{}, b interface{}) {})

	var sg mocks.StatisticsGateway

	c.Scan("statistics_gateway", &sg)
	stats := sg.CollectStatistics()

	booksCount := stats["books_count"]
	authorsCount := stats["authors_count"]

	if booksCount != 2 {
		t.Errorf("Wrong books count provided '%d', expected count is '%d'", booksCount, 2)
	}

	if authorsCount != 3 {
		t.Errorf("Wrong authors count provided '%d', expected count is '%d'", authorsCount, 3)
	}
}

func TestEventMerging(t *testing.T) {
	evCont1 := NewEventsContainer()
	evCont1.addDependencyObserver("event1", "service1", func(a, b interface{}) {})
	evCont1.registerDependencyEvent("event1", "dependency1")

	evCont2 := NewEventsContainer()

	funcToGetNotificationIsCalled := false
	funcToGetNotificationAboutDependency := func(service, dependency interface{}) {
		funcToGetNotificationIsCalled = true
		if service.(string) != "someServiceInstance" {
			t.Errorf(
				"Service '%s' with id '%s'rather than serivce '%s' should get a dependency notification grouped by Event '%s' with dependency '%s' as context",
				"someServiceInstance",
				service.(string),
				"observerId2",
				"event2",
				"dependency2",
			)
		}
		if dependency.(string) != "someDependencyInstance" {
			t.Errorf(
				"A Dependency notification '%s' is expected rather than rather than '%s'",
				"someDependencyInstance",
				dependency.(string),
			)
		}
	}
	evCont2.addDependencyObserver("event2", "observerId2", funcToGetNotificationAboutDependency)
	evCont2.registerDependencyEvent("event2", "dependency2")

	evCont1.merge(*evCont2)

	dependencyInstance := "someDependencyInstance"
	containerMock := ContainerInterfaceMock{service: dependencyInstance}
	serviceInstance := "someServiceInstance"

	evCont1.collectDependencyEventsForService(&containerMock, "observerId2", serviceInstance)

	if !funcToGetNotificationIsCalled {
		t.Errorf(
			"Service '%s' with id '%s' should get a dependency notification grouped by Event '%s' with dependency '%s' as context but none was received",
			"someServiceInstance",
			"observerId2",
			"event2",
			"dependency2",
		)
	}
}

func TestNoConstructorCalledIfItsNotInvolvedInDependencyChain(t *testing.T) {
	runtimeContainer := CreateContainer()
	failingNewFunc := func() string {
		panic("Should not be called")
	}
	runtimeContainer.AddNewMethod("some_failing_service", failingNewFunc)
	runtimeContainer.RegisterDependencyEvent("statistics_provider", "some_failing_service")

	runtimeContainer.Get("book_shelve", true)
}
