package tests

import (
	"testing"
	"github.com/breathbath/gotainer/examples"
)

func TestAllPresetDefinitions(t *testing.T) {
	cont := examples.CreateContainer()

	var bc examples.BookCreator
	cont.GetTypedService("book_creator", &bc)

	var bs examples.BookStorage
	cont.GetTypedService("book_storage", &bs)

	var finder examples.BookFinder
	cont.GetTypedService("book_finder", &finder)

	book, result := finder.FindBook("one")
	if !result || book.Title != "FirstBook" || book.Author != "FirstAuthor" || book.Id != "One" {
		t.Errorf(
			"Book finder cannot find a book, probably its intialised correctly",
		)
	}

	var staticFinder examples.BookFinder
	cont.GetTypedService("book_finder_declared_statically", &staticFinder)

	book, result = finder.FindBook("two")
	if !result || book.Title != "SecondBook" || book.Author != "FirstAuthor" || book.Id != "Two" {
		t.Errorf(
			"Book finder cannot find a book, probably it's intialised incorrectly",
		)
	}
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