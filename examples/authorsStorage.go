package examples

type AuthorsStorage struct {
	db FakeDb
}

func NewAuthorsStorage(db FakeDb) AuthorsStorage {
	return AuthorsStorage{db}
}

func (bs AuthorsStorage) GetStatistics() (string, int) {
	return "authors_count", bs.db.CountItems("authors")
}