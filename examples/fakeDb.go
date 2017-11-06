package examples

import "fmt"

type FakeDb struct {
	data             map[string]map[string]string
	connectionString string
}

func NewFakeDb(connectionString string) FakeDb {
	//to make sure that a correct value is provided from the config
	if connectionString != "someConnectionString" {
		panic(
			fmt.Errorf(
				"Incorrect connection string provided '%s', expected connection string: %s",
				connectionString,
				"someConnectionString",
			))
	}
	return FakeDb{connectionString: connectionString, data: map[string]map[string]string{
		"books": {
			"one": "One;FirstBook;FirstAuthor",
			"two": "Two;SecondBook;FirstAuthor",
		},
	}}
}

func (fdb FakeDb) FindInTable(tableName, id string) (string, bool) {
	tableData, bookName := map[string]string{}, ""
	found := false
	tableData, found = fdb.data[tableName]
	if found {
		bookName, found = tableData[id]
	}

	return bookName, found
}
