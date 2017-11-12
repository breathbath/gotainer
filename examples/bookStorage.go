package examples

//BookStorage repository for books data
type BookStorage struct {
	db FakeDb
}

//NewBookStorage BookStorage constructor
func NewBookStorage(db FakeDb) BookStorage {
	return BookStorage{db}
}

//FindBookData get a book by id
func (bs BookStorage) FindBookData(id string) (string, bool) {
	return bs.db.FindInTable("books", id)
}

//GetStatistics returns books count
func (bs BookStorage) GetStatistics() (string, int) {
	return "books_count", bs.db.CountItems("books")
}

//GetAllBooks get some collection of books
func GetAllBooks() []Book {
	return []Book{Book{"1", "Book1", "Author1"}, Book{"2", "Book2", "Author2"}}
}
