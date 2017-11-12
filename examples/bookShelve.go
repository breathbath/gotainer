package examples

//BookShelve holds a temp collection of books
type BookShelve struct {
	books []Book
}

//Main constructor
func NewBookShelve() *BookShelve {
	return &BookShelve{[]Book{}}
}

//Adding a book to collection
func (bs *BookShelve) Add(book Book) {
	bs.books = append(bs.books, book)
}

//Returns all books
func (bs *BookShelve) GetBooks() []Book {
	return bs.books
}
