package examples

type BookShelve struct {
	books []Book
}

func NewBookShelve() *BookShelve {
	return &BookShelve{[]Book{}}
}

func (bs *BookShelve) Add(book Book) {
	bs.books = append(bs.books, book)
}

func (bs *BookShelve) GetBooks() [] Book {
	return bs.books
}
