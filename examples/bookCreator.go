package examples

import "strings"

type BookCreator struct{}

func (bc BookCreator) CreateBook(bookData string) Book {
	bookFields := strings.Split(bookData, ";")
	return Book{Id: bookFields[0], Title: bookFields[1], Author: bookFields[2]}
}
