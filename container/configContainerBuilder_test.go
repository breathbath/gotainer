package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestAllDependencyTypesCreatedFromConfigCorrectly(t *testing.T) {
	container, err := buildContainerFromMockedConfig()
	if err != nil {
		t.Error(err)
		return
	}

	var config mocks.Config
	container.Scan("config", &config)
	if config.GetValue("staticFilesUrl") != "http://static.me/" {
		t.Error("Wrong 'config' service is returned from the container")
	}

	if container.Get("connection_string", true).(string) != "someConnectionString" {
		t.Error("Wrong 'connection_string' is returned from the container")
	}

	var fakeDb mocks.FakeDb
	container.Scan("db", &fakeDb)
	if fakeDb.CountItems("books") != 2 {
		t.Error("Wrong 'db' is returned from the container")
	}

	var statsGateway mocks.StatisticsGateway
	container.Scan("statistics_gateway", &statsGateway)
	stats := statsGateway.CollectStatistics()
	if stats["authors_count"] != 3 {
		t.Error("A wrongly initialised 'statistics_gateway' is returned from the container")
	}
}

func TestParameters(t *testing.T) {
	container, err := buildContainerFromMockedConfig()
	if err != nil {
		t.Error(err)
		return
	}

	AssertExpectedDependency(container, "param1", "value1", t)
	AssertExpectedDependency(container, "param2", 123, t)
	AssertExpectedDependency(container, "EnableLogging", true, t)

	var logger mocks.InMemoryLogger
	container.Scan("logger", &logger)
	logger.Log("message")

	expectedResult := "message"
	providedResult := logger.GetMessages()[0]
	if providedResult != expectedResult {
		t.Errorf(
			"Unexpected logged result %s, expected result is %s",
			providedResult,
			expectedResult,
		)
	}
}

func TestBuildingWithWrongObserverDeclaration(t *testing.T) {
	config := Tree{
		Node{
			NewFunc: mocks.NewBookShelve,
			ID:      "book_shelve",
		},
		Node{
			Ob: Observer{
				"some_event",
				"some_name",
				"not_function",
			},
		},
	}

	_, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(config)
	AssertError(
		err,
		"A function is expected rather than 'string' [check 'Node: {ID: ; ServiceNames: []; Event: {Name: ; Service: ;}; Observer: {Name: some_name; Event: some_event;}}' service]",
		t,
	)
}

func TestConfigMerge(t *testing.T) {
	configTree := getMockedConfigTree()
	configTreeToMerge := Tree{
		Node{
			NewFunc: mocks.NewBookShelve,
			ID:      "book_shelve",
		},
		Node{
			NewFunc:      mocks.NewBookRevision,
			ID:           "book_revision",
			ServiceNames: Services{"book_finder_declared_statically"},
		},
	}

	container, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(configTree, configTreeToMerge)
	if err != nil {
		t.Error(err)
		return
	}

	var bookShelve mocks.BookShelve
	container.Scan("book_shelve", &bookShelve)
	bookShelve.Add(mocks.Book{Id: "someBook"})
	book := bookShelve.GetBooks()[0]

	if book.Id != "someBook" {
		t.Error("A wrongly working 'book_shelve' is returned from the container after config merge")
	}

	var connectionString string
	container.Scan("connection_string", &connectionString)

	if connectionString != "someConnectionString" {
		t.Error("A wrong service declaration for 'connection_string' is returned from the container after config merge")
	}
}

func TestConfigServiceExistence(t *testing.T) {
	configTree := getMockedConfigTree()
	if !configTree.ServiceExists("connection_string") {
		t.Errorf("It is expected that 'connection_string' parameter exits but it's not")
	}

	if configTree.ServiceExists("some_non_existing_service") {
		t.Errorf("It is expected that 'some_non_existing_service' parameter does not exist, but it does")
	}
}

func TestConfigGarbageCollectionSuccess(t *testing.T) {
	container, err := buildContainerFromMockedConfig()
	if err != nil {
		t.Errorf("Unexpected error %v during the garbage collection", err)
		return
	}

	fakeDb := container.Get("db", true).(*mocks.FakeDb)

	if fakeDb.WasDestroyed() {
		t.Error("FakeDb should not have been destroyed before the garbage collect call")
		return
	}

	err = container.CollectGarbage()
	if err != nil {
		t.Errorf("Unexpected error %v during the garbage collection", err)
		return
	}

	if !fakeDb.WasDestroyed() {
		t.Error("FakeDb should have been destroyed after the garbage collect call")
	}
}
