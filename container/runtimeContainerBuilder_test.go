package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
)

func TestAllDependencyTypesCreatedFromConfigCorrectly(t *testing.T) {
	configTree := getMockedConfigTree()
	container := RuntimeContainerBuilder{}.BuildContainer(configTree)

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
