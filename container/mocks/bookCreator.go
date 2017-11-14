package mocks

import "strings"

//BookCreator converts string ; separated data into Book entity
type BookCreator struct{}

//CreateBook converts string into a Book entity
func (bc BookCreator) CreateBook(bookData string) Book {
	bookFields := strings.Split(bookData, ";")
	return Book{Id: bookFields[0], Title: bookFields[1], Author: bookFields[2]}
}
