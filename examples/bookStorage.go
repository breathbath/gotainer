package examples

type BookStorage struct {
	db FakeDb
}

func NewBookStorage(db FakeDb) BookStorage {
	return BookStorage{db}
}

func (bs BookStorage) FindBookData(id string) (string, bool) {
	return bs.db.FindInTable("books", id)
}

func (bs BookStorage) GetStatistics() (string, int) {
	return "books_count", bs.db.CountItems("books")
}