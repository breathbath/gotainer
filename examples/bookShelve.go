package examples

//BookShelve holds a temp collection of books
type BookShelve struct {
	books []Book
}

//NewBookShelve BookShelve constructor
func NewBookShelve() *BookShelve {
	return &BookShelve{[]Book{}}
}

//Add a book to collection
func (bs *BookShelve) Add(book Book) {
	bs.books = append(bs.books, book)
}

//GetBooks returns all books
func (bs *BookShelve) GetBooks() []Book {
	return bs.books
}
