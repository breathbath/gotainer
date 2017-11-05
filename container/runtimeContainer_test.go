package container

import (
	"strings"
	"testing"
)

type Book struct {
	Id     string
	Title  string
	Author string
}

type StorageService struct {
	booksTable map[string]string
}

func NewStorageService() StorageService {
	table := map[string]string{
		"one": "One;FirstBook;FirstAuthor",
		"two": "Two;SecondBook;FirstAuthor",
	}
	return StorageService{table}
}

func (ss StorageService) FindBookData(id string) (string, bool) {
	bookName, ok := ss.booksTable[id]

	return bookName, ok
}

type BookCreator struct{}

func (bc BookCreator) CreateBook(bookData string) Book {
	bookFields := strings.Split(bookData, ";")
	return Book{Id: bookFields[0], Title: bookFields[1], Author: bookFields[2]}
}

type BookFinder struct {
	Storage     StorageService
	BookCreator BookCreator
}

func NewBookFinder(storage StorageService, bookCreator BookCreator) BookFinder {
	return BookFinder{storage, bookCreator}
}

func (bc BookFinder) FindBook(id string) (Book, bool) {
	bookData, found := bc.Storage.FindBookData(id)
	if !found {
		return Book{}, false
	}

	book := bc.BookCreator.CreateBook(bookData)

	return book, true
}

func CreateContainer() Container {
	container := NewRuntimeContainer()

	container.AddNoArgumentsConstructor("book_creator", func() (interface{}, error) {
		return BookCreator{}, nil
	})

	container.AddNoArgumentsConstructor("book_storage", func() (interface{}, error) {
		return NewStorageService(), nil
	})

	container.AddConstructor("book_finder", func(c Container) (interface{}, error) {
		var bc BookCreator
		c.GetTypedService("book_creator", &bc)

		var bs StorageService
		c.GetTypedService("book_storage", &bs)

		return NewBookFinder(bs, bc), nil
	})

	container.AddTypedConstructor("book_finder_declared_statically", NewBookFinder, "book_storage", "book_creator")

	return container
}

func TestAllPresetDefinitions(t *testing.T) {
	container := CreateContainer()

	var bc BookCreator
	container.GetTypedService("book_creator", &bc)

	var bs StorageService
	container.GetTypedService("book_storage", &bs)

	var finder BookFinder
	container.GetTypedService("book_finder", &finder)

	book, result := finder.FindBook("one")
	if !result || book.Title != "FirstBook" || book.Author != "FirstAuthor" || book.Id != "One" {
		t.Errorf(
			"Book finder cannot find a book, probably its intialised correctly",
		)
	}

	var staticFinder BookFinder
	container.GetTypedService("book_finder_declared_statically", &staticFinder)

	book, result = finder.FindBook("two")
	if !result || book.Title != "SecondBook" || book.Author != "FirstAuthor" || book.Id != "Two" {
		t.Errorf(
			"Book finder cannot find a book, probably its intialised correctly",
		)
	}
}

func TestWrongMixedDependenciesInStaticCall(t *testing.T) {
	defer ExpectPanic("Cannot use the provided service 'book_creator' of type 'container.BookCreator' as 'container.StorageService' in the constructor function call", t)

	container := CreateContainer()
	container.AddTypedConstructor("anotherFinder", NewBookFinder, "book_creator", "book_storage")

	var finder BookFinder
	container.GetTypedService("anotherFinder", &finder)
}

func TestWrongDependencyRequested(t *testing.T) {
	defer ExpectPanic("Unknown service 'lala'", t)
	container := CreateContainer()

	container.GetTypedService("lala", nil)
}
