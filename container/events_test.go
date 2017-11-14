package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
)

func TestContainerEvents(t *testing.T) {
	c := CreateContainer()
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
