package container

import (
	"testing"
	"fmt"
	"github.com/breathbath/gotainer/container/mocks"
)

func setupTwoContainers() (*RuntimeContainer, *RuntimeContainer) {
	container1 := NewRuntimeContainer()
	container1.AddNewMethod("book_shelve", mocks.NewBookShelve)

	bookShelve := container1.Get("book_shelve", true).(*mocks.BookShelve)
	bookShelve.Add(mocks.Book{Id: "1", Title: "Book"})

	container1.AddConstructor("memory_cache", func(c Container) (interface{}, error) {
		return mocks.NewInMemoryCache(), nil
	})

	container2 := NewRuntimeContainer()
	container2.AddNewMethod("config", mocks.NewConfig)

	return container1, container2
}

//TestContainerMerge expects services from both containers to be available
func TestContainerMerge(t *testing.T) {
	container1, container2 := setupTwoContainers()
	container2.Merge(container1)

	var config mocks.Config
	container2.Scan("config", &config)

	if config.GetValue("fakeDbConnectionString") != "someConnectionString" {
		t.Error("'Config' service merged into the container has a wrong definition")
	}

	bookShelve := container2.Get("book_shelve", true).(*mocks.BookShelve)

	assertBookShelveContainsBook(bookShelve, 1, "1", t)
}

//TestNonUniqueServicesMerge handles the test case when both containers have a service with the same name
func TestNonUniqueServicesMerge(t *testing.T) {

	container1, container2 := setupTwoContainers()
	container2.AddNewMethod("book_shelve", mocks.NewBookShelve)
	defer ExpectPanic(fmt.Sprintf("Cannot merge containers because of non unique service id '%s'", "book_shelve"), t)
	container1.Merge(container2)
}

func assertBookShelveContainsBook(bookShelve *mocks.BookShelve, expectedBooksAmount int, expectedBookId string, t *testing.T) {
	books := bookShelve.GetBooks()
	if expectedBooksAmount == 0 && len(books) != 0 {
		t.Error("The service 'book_shelve' shouldn't contain books but it has at least one book inside")
		return;
	}

	if len(books) != expectedBooksAmount || books[0].Id != expectedBookId {
		t.Error("The service 'book_shelve' should contain a book added before then container merge")
	}
}
