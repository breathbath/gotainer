package mocks

//AuthorsStorage repository for authors data
type AuthorsStorage struct {
	db *FakeDb
}

//NewAuthorsStorage constructor
func NewAuthorsStorage(db *FakeDb) AuthorsStorage {
	return AuthorsStorage{db}
}

//GetStatistics counts authors
func (bs AuthorsStorage) GetStatistics() (string, int) {
	return "authors_count", bs.db.CountItems("authors")
}
