package tests

import (
	"github.com/breathbath/gotainer/examples"
	"testing"
)

func TestContainerCache(t *testing.T) {
	c := examples.CreateContainer()

	var bookShelve examples.BookShelve

	c.Scan("book_shelve", &bookShelve)

	bookShelve.Add(examples.Book{Id: "123", Title: "Book1", Author: "Author1"})

	var bookShelve2 examples.BookShelve

	c.Scan("book_shelve", &bookShelve2)
	AssertBookShelveHasBook("123", "Book1", "Author1", t, bookShelve)

	c.ScanNonCached("book_shelve", &bookShelve2)
	if len(bookShelve2.GetBooks()) != 0 {
		t.Errorf("Book shelve should be empty if container is asked for uncached dependency")
	}
}

func AssertBookShelveHasBook(id, title, author string, t *testing.T, bs examples.BookShelve) {
	if len(bs.GetBooks()) != 1 {
		t.Errorf("Book shelve should contain books after adding")
		return
	}

	book := bs.GetBooks()[0]

	if book.Id != id || book.Title != title || book.Author != author {
		t.Errorf("Wrong book returned from shelve")
	}
}
