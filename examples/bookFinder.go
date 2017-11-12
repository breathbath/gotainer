package examples

//BookFinder searches books
type BookFinder struct {
	Storage     BookStorage
	BookCreator BookCreator
}

//NewBookFinder constructor for BookFinder
func NewBookFinder(storage BookStorage, bookCreator BookCreator) BookFinder {
	return BookFinder{storage, bookCreator}
}

//FindBook get a book by id
func (bc BookFinder) FindBook(id string) (Book, bool) {
	bookData, found := bc.Storage.FindBookData(id)
	if !found {
		return Book{}, false
	}

	book := bc.BookCreator.CreateBook(bookData)

	return book, true
}
