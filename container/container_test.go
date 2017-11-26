package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
	_ "reflect"
	_ "fmt"
	"reflect"
)

func TestAllPresetDefinitions(t *testing.T) {
	cont := CreateContainer()

	var bc mocks.BookCreator
	cont.Scan("book_creator", &bc)

	var bs mocks.BookStorage
	cont.Scan("book_storage", &bs)

	var finder mocks.BookFinder
	cont.Scan("book_finder", &finder)

	book, result := finder.FindBook("one")
	if !result || book.Title != "FirstBook" || book.Author != "FirstAuthor" || book.Id != "One" {
		t.Errorf(
			"Book finder cannot find a book, probably its intialised correctly",
		)
	}

	var staticFinder mocks.BookFinder
	cont.Scan("book_finder_declared_statically", &staticFinder)

	book, result = finder.FindBook("two")
	if !result || book.Title != "SecondBook" || book.Author != "FirstAuthor" || book.Id != "Two" {
		t.Errorf(
			"Book finder cannot find a book, probably it's intialised incorrectly",
		)
	}
}

func TestStaticParameterDependency(t *testing.T) {
	container := CreateContainer()
	var bookLinkProvider mocks.BookLinkProvider
	container.Scan("book_link_provider", &bookLinkProvider)

	url := bookLinkProvider.GetLink("one")
	expectedUrl := "http://static.me/FirstBook"
	if url != expectedUrl {
		t.Errorf(
			"Unexpected book url '%s' fetched, expected url is '%s'.",
			url,
			expectedUrl,
		)
	}
}

func TestPointerAndInterfaceDependencies(t *testing.T) {
	var bookDownloader mocks.BookDownloader
	container := CreateContainer()
	container.Scan("book_downloader", &bookDownloader)

	fetchedContent := bookDownloader.DownloadBook("two")
	expectedFetchedContentd := "Fetched from http://static.me/SecondBook"
	if fetchedContent != expectedFetchedContentd {
		t.Errorf(
			"Unexpected book content '%s' fetched, expected content is '%s'.",
			fetchedContent,
			expectedFetchedContentd,
		)
	}
}

func TestNewMethodWithTwoReturns(t *testing.T) {
	newMethodWithTwoReturns1 := func() (mocks.Book, error) {
		return mocks.Book{Id: "123"}, nil
	}

	newMethodWithTwoReturns2 := func() (error, mocks.Book) {
		return nil, mocks.Book{Id: "456"}
	}

	container := CreateContainer()
	container.AddNewMethod("some_book1", newMethodWithTwoReturns1)
	container.AddNewMethod("some_book2", newMethodWithTwoReturns2)

	assertExpectedBookInContainer(container, "some_book1", "123", t)
	assertExpectedBookInContainer(container, "some_book2", "456", t)
}

func assertExpectedBookInContainer(container Container, bookServiceId, expectedBookId string, t *testing.T) {
	book := container.Get(bookServiceId, true).(mocks.Book)

	if book.Id != expectedBookId {
		t.Errorf(
			"New method with 2 returns should return a book with id '%s', but none was returned",
			expectedBookId,
		)
	}
}

type A struct {
	A1 int
}

func TestPointerDependency(t *testing.T) {
	c := CreateContainer()
	c.AddNewMethod("some_pointer_service", mocks.NewBookShelve)

	//var a *A
	//b := A{1}
	//va := reflect.ValueOf(&a).Elem()
	//v := reflect.New(va.Type().Elem())
	//vb := reflect.ValueOf(&b).Elem()
	//vbIndirect := reflect.Indirect(vb)
	//vba := vbIndirect.Convert(reflect.TypeOf(a).Elem())
	//va.Set(v)
	//
	//
	//fmt.Println(vba)
	//return
	//	var bookShelve *mocks.BookShelve
	//
	//	intType := reflect.TypeOf(bookShelve).Elem()
	//	intPtr2 := reflect.New(intType)
	//	//book := mocks.Book{Id:"id"}
	//	//bookShelve.Add(book)
	//
	//	//bbb := intPtr2.Elem().Interface().(mocks.BookShelve)
	//return

	//var a *string
	//c.Scan("static_files_url", a)
	//fmt.Println(a)
	//return

	var bookShelve mocks.BookShelve
	conv(&bookShelve)
	//va := reflect.ValueOf(&bookShelve).Elem()
	//v := reflect.New(va.Type().Elem())
	//va.Set(v)
	c.Scan("some_pointer_service", &bookShelve)

	//var bookShelve *mocks.BookShelve
	//c.Scan("some_pointer_service", bookShelve)
	book := mocks.Book{Id: "id"}
	bookShelve.Add(book)

	bookReturned := bookShelve.GetBooks()[0]
	if bookReturned.Id != "id" {
		t.Errorf("Returned pointer object has a wrong type")
	}
}

func conv(lal *mocks.BookShelve) {
	va := reflect.ValueOf(&lal).Elem()
	v := reflect.New(va.Type().Elem())
	va.Set(v)
}

func TestCheckNotFails(t *testing.T) {
	cont := CreateContainer()
	cont.Check()
}
