package tests

import (
	"testing"
	"github.com/breathbath/gotainer/container"
	"github.com/breathbath/gotainer/examples"
	"fmt"
)

func setupTwoContainers() (*container.RuntimeContainer, *container.RuntimeContainer) {
	container1 := container.NewRuntimeContainer()
	container1.AddNewMethod("book_shelve", examples.NewBookShelve)

	bookShelve := container1.Get("book_shelve", true).(*examples.BookShelve)
	bookShelve.Add(examples.Book{Id: "1", Title: "Book"})

	container1.AddConstructor("memory_cache", func(c container.Container) (interface{}, error) {
		return examples.NewInMemoryCache(), nil
	})

	container2 := container.NewRuntimeContainer()
	container2.AddNewMethod("config", examples.NewConfig)

	return container1, container2
}

//TestContainerMerge expects services from both containers to be available
func TestContainerMerge(t *testing.T) {
	container1, container2 := setupTwoContainers()
	container1.Merge(container2)

	var config examples.Config
	container1.Scan("config", &config)

	if config.GetValue("fakeDbConnectionString") != "someConnectionString" {
		t.Error("'Config' service merged into the container has a wrong definition")
	}

	bookShelve := container1.Get("book_shelve", true).(*examples.BookShelve)

	assertBookShelveContainsBook(bookShelve, 1, "1", t)
}

//TestNonUniqueServicesMerge handles the test case when both containers have a service with the same name
func TestNonUniqueServicesMerge(t *testing.T) {

	container1, container2 := setupTwoContainers()
	container2.AddNewMethod("book_shelve", examples.NewBookShelve)
	defer ExpectPanic(fmt.Sprintf("Cannot merge containers because of non unique service id '%s'", "book_shelve"), t)
	container1.Merge(container2)
}

func assertBookShelveContainsBook(bookShelve *examples.BookShelve, expectedBooksAmount int, expectedBookId string, t *testing.T) {
	books := bookShelve.GetBooks()
	if expectedBooksAmount == 0 && len(books) != 0 {
		t.Error("The service 'book_shelve' shouldn't contain books but it has at least one book inside")
		return;
	}

	if len(books) != expectedBooksAmount || books[0].Id != expectedBookId {
		t.Error("The service 'book_shelve' should contain a book added before then container merge")
	}
}
