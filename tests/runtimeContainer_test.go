package tests

import (
	"testing"
	"github.com/breathbath/gotainer/examples"
)

func TestAllPresetDefinitions(t *testing.T) {
	container := examples.CreateContainer()

	var bc examples.BookCreator
	container.GetTypedService("book_creator", &bc)

	var bs examples.BookStorage
	container.GetTypedService("book_storage", &bs)

	var finder examples.BookFinder
	container.GetTypedService("book_finder", &finder)

	book, result := finder.FindBook("one")
	if !result || book.Title != "FirstBook" || book.Author != "FirstAuthor" || book.Id != "One" {
		t.Errorf(
			"Book finder cannot find a book, probably its intialised correctly",
		)
	}

	var staticFinder examples.BookFinder
	container.GetTypedService("book_finder_declared_statically", &staticFinder)

	book, result = finder.FindBook("two")
	if !result || book.Title != "SecondBook" || book.Author != "FirstAuthor" || book.Id != "Two" {
		t.Errorf(
			"Book finder cannot find a book, probably it's intialised incorrectly",
		)
	}
}

func TestWrongMixedDependenciesInStaticCall(t *testing.T) {
	defer ExpectPanic("Cannot use the provided service 'book_creator' of type 'examples.BookCreator' as 'examples.BookStorage' in the constructor function call", t)

	container := examples.CreateContainer()
	container.AddTypedConstructor("anotherFinder", examples.NewBookFinder, "book_creator", "book_storage")

	var finder examples.BookFinder
	container.GetTypedService("anotherFinder", &finder)
}

func TestWrongDependencyRequested(t *testing.T) {
	defer ExpectPanic("Unknown service 'lala'", t)
	container := examples.CreateContainer()

	container.GetTypedService("lala", nil)
}

func TestStaticParameterDependency(t *testing.T) {
	container := examples.CreateContainer()
	var bookLinkProvider examples.BookLinkProvider
	container.GetTypedService("book_link_provider", &bookLinkProvider)

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
	var bookDownloader examples.BookDownloader
	container := examples.CreateContainer()
	container.GetTypedService("book_downloader", &bookDownloader)

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

func TestIncompatibleInterfaces(t *testing.T) {
	defer ExpectPanic("Cannot use the provided service 'incompatible_cache' of type 'examples.IncompatibleCache' as 'examples.Cache' in the constructor function call", t)
	container := examples.CreateContainer()
	container.AddTypedConstructor("incompatible_cache", examples.NewIncompatibleCache)
	container.AddTypedConstructor(
		"wrong_downloader",
		examples.NewBookDownloader,
		"incompatible_cache",
		"book_link_provider",
		"book_finder",
		"web_fetcher",
	)

	container.GetTypedService("incompatible_cache", nil)
}

func TestCheck(t *testing.T) {
	defer ExpectPanic("Cannot convert created value of type 'int' to expected destination value 'BookCreator' for createdDependency declaration wrong_book_creator", t)
	container := examples.CreateContainer()
	container.Check()
}
