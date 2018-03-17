package mocks

import (
	"fmt"
)

//BookRevision service to check book status
type BookRevision struct {
	bf BookFinder
}

//NewBookRevision BookRevision constructor
func NewBookRevision(bf BookFinder) BookRevision {
	return BookRevision{bf}
}

//IsBookDataComplete checks the completion of book data
func (br BookRevision) IsBookDataComplete(id string) (bool, error) {
	bookData, found := br.bf.FindBook(id)
	if !found {
		return false, fmt.Errorf("Cannot check the completion of a non-existing book '%s'", id)
	}

	return bookData.Id != "" && bookData.Title != "" && bookData.Author != "", nil
}
