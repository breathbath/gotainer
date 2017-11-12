package examples

import "fmt"

//FakeDb mock for simulating a db engine
type FakeDb struct {
	data             map[string]map[string]string
	connectionString string
}

//FakeDb main constructor
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
		"authors": {
			"1a": "FirstAuthor",
			"2a": "SecondAuthor",
			"3a": "ThirdAuthor",
		},
	}}
}

//FindInTable simulates a search functionality in a db table
func (fdb FakeDb) FindInTable(tableName, id string) (string, bool) {
	tableData, bookName := map[string]string{}, ""
	found := false
	tableData, found = fdb.data[tableName]
	if found {
		bookName, found = tableData[id]
	}

	return bookName, found
}

//CountItems counts items in a simulated table
func (fdb FakeDb) CountItems(tableName string) int {
	tableData, found := fdb.data[tableName]
	if found {
		return len(tableData)
	}
	return 0
}
