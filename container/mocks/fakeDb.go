package mocks

import "fmt"

//FakeDb mock for simulating a db engine
type FakeDb struct {
	data             map[string]map[string]string
	connectionString string
	isDestroyed      bool
}

//NewFakeDb main constructor
func NewFakeDb(connectionString string) *FakeDb {
	//to make sure that a correct value is provided from the config
	if connectionString != "someConnectionString" {
		panic(
			fmt.Errorf(
				"Incorrect connection string provided '%s', expected connection string: %s",
				connectionString,
				"someConnectionString",
			))
	}
	return &FakeDb{
		connectionString: connectionString,
		data: map[string]map[string]string{
			"books": {
				"one": "One;FirstBook;FirstAuthor",
				"two": "Two;SecondBook;FirstAuthor",
			},
			"authors": {
				"1a": "FirstAuthor",
				"2a": "SecondAuthor",
				"3a": "ThirdAuthor",
			},
		},
		isDestroyed: false,
	}
}

//FindInTable simulates a search functionality in a db table
func (fdb *FakeDb) FindInTable(tableName, id string) (string, bool) {
	if fdb.isDestroyed {
		panic("Database was already destroyed")
	}
	var (
		tableData map[string]string
		bookName string
		found bool
	)

	tableData, found = fdb.data[tableName]
	if found {
		bookName, found = tableData[id]
	}

	return bookName, found
}

//Destroy garbage collection func
func (fdb *FakeDb) Destroy() error {
	fdb.isDestroyed = true
	return nil
}

//WasDestroyed mocked func to make sure if the Destroy method was called before
func (fdb FakeDb) WasDestroyed() bool {
	return fdb.isDestroyed
}

//CountItems counts items in a simulated table
func (fdb *FakeDb) CountItems(tableName string) int {
	if fdb.isDestroyed {
		panic("Database was already destroyed")
	}

	tableData, found := fdb.data[tableName]
	if found {
		return len(tableData)
	}
	return 0
}
