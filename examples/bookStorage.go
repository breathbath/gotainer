package examples

//BookStorage repository for books data
type BookStorage struct {
	db FakeDb
}

//Main contructor
func NewBookStorage(db FakeDb) BookStorage {
	return BookStorage{db}
}

//Finds a book by id
func (bs BookStorage) FindBookData(id string) (string, bool) {
	return bs.db.FindInTable("books", id)
}

//Returns books count
func (bs BookStorage) GetStatistics() (string, int) {
	return "books_count", bs.db.CountItems("books")
}

//Returns some collection of books
func GetAllBooks() []Book {
	return []Book{Book{"1", "Book1", "Author1"}, Book{"2", "Book2", "Author2"}}
}
