package examples

//BookFinder searches books
type BookFinder struct {
	Storage     BookStorage
	BookCreator BookCreator
}

//Constructor
func NewBookFinder(storage BookStorage, bookCreator BookCreator) BookFinder {
	return BookFinder{storage, bookCreator}
}

//Finds a book by id
func (bc BookFinder) FindBook(id string) (Book, bool) {
	bookData, found := bc.Storage.FindBookData(id)
	if !found {
		return Book{}, false
	}

	book := bc.BookCreator.CreateBook(bookData)

	return book, true
}
