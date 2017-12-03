package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
)

func TestAllDependencyTypesCreatedFromConfigCorrectly(t *testing.T) {
	configTree := getMockedConfigTree()
	container := RuntimeContainerBuilder{}.BuildContainerFromConfig(configTree)

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
	configTree := getMockedConfigTree()
	container := RuntimeContainerBuilder{}.BuildContainerFromConfig(configTree)

	AssertExpectedDependency(container, "param1", "value1", t)
	AssertExpectedDependency(container, "param2", 123, t)
	AssertExpectedDependency(container, "EnableLogging", true, t)

	var logger mocks.InMemoryLogger
	container.Scan("logger", &logger)
	logger.Log("message")

	expectedResult :=  "message"
	providedResult := logger.GetMessages()[0]
	if providedResult != expectedResult {
		t.Errorf(
			"Unexpected logged result %s, expected result is %s",
			providedResult,
			expectedResult,
		)
	}
}

func TestConfigMerge(t *testing.T) {
	configTree := getMockedConfigTree()
	configTreeToMerge := Tree{
		Node{
			NewFunc: mocks.NewBookShelve,
			Id:      "book_shelve",
		},
	}

	container := RuntimeContainerBuilder{}.BuildContainerFromConfig(configTree, configTreeToMerge)

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
