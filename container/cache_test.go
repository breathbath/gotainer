package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestContainerCache(t *testing.T) {
	c := CreateContainer()

	var bookShelve mocks.BookShelve

	c.Scan("book_shelve", &bookShelve)

	bookShelve.Add(mocks.Book{Id: "123", Title: "Book1", Author: "Author1"})

	var bookShelve2 mocks.BookShelve

	c.Scan("book_shelve", &bookShelve2)
	AssertBookShelveHasBook("123", "Book1", "Author1", t, bookShelve)

	c.ScanNonCached("book_shelve", &bookShelve2)
	if len(bookShelve2.GetBooks()) != 0 {
		t.Errorf("Book shelve should be empty if container is asked for uncached dependency")
	}
}

func TestNonCachedRequestForMismatchedDestinationType(t *testing.T) {
	c := CreateContainer()
	defer ExpectPanic("Cannot convert created value of type 'Config' to expected destination value 'Book' for createdDependency declaration config [check 'config' service]", t)
	var book mocks.Book

	c.ScanNonCached("config", &book)
}

func AssertBookShelveHasBook(id, title, author string, t *testing.T, bs mocks.BookShelve) {
	if len(bs.GetBooks()) != 1 {
		t.Errorf("Book shelve should contain books after adding")
		return
	}

	book := bs.GetBooks()[0]

	if book.Id != id || book.Title != title || book.Author != author {
		t.Errorf("Wrong book returned from shelve")
	}
}
