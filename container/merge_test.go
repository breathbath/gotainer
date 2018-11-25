package container

import (
	"fmt"
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
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
	container2.AddNewMethod("Config", mocks.NewConfig)

	return container1, container2
}

//TestContainerMerge expects Services from both containers to be available
func TestContainerMerge(t *testing.T) {
	container1, container2 := setupTwoContainers()
	container2.Merge(container1)

	var config mocks.Config
	container2.Scan("Config", &config)

	if config.GetValue("fakeDbConnectionString") != "someConnectionString" {
		t.Error("'Config' Service merged into the container has a wrong definition")
	}

	bookShelve := container2.Get("book_shelve", true).(*mocks.BookShelve)

	assertBookShelveContainsBook(bookShelve, 1, "1", t)
}

//TestNonUniqueServicesMerge handles the test case when both containers have a Service with the same Name
func TestNonUniqueServicesMerge(t *testing.T) {
	container1, container2 := setupTwoContainers()
	container2.AddNewMethod("book_shelve", mocks.NewBookShelve)
	defer ExpectPanic(t, fmt.Sprintf("Cannot merge containers because of non unique Service id '%s'", "book_shelve"))
	container1.Merge(container2)
}

func TestConflictingMerge(t *testing.T) {
	defer ExpectPanic(t, "Cannot merge containers because of non unique Service id 'serviceA'")
	cont1 := NewRuntimeContainer()
	cont1.AddConstructor("serviceA", func(c Container) (interface{}, error) {
		return "serviceA", nil
	})

	cont2 := NewRuntimeContainer()
	cont2.AddConstructor("serviceA", func(c Container) (interface{}, error) {
		return "serviceA", nil
	})

	cont1.Merge(cont2)
}

func assertBookShelveContainsBook(bookShelve *mocks.BookShelve, expectedBooksAmount int, expectedBookId string, t *testing.T) {
	books := bookShelve.GetBooks()
	if expectedBooksAmount == 0 && len(books) != 0 {
		t.Error("The Service 'book_shelve' shouldn't contain books but it has at least one book inside")
		return
	}

	if len(books) != expectedBooksAmount || books[0].Id != expectedBookId {
		t.Error("The Service 'book_shelve' should contain a book added before then container merge")
	}
}
