package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
)

func TestWrongMixedDependenciesInStaticCall(t *testing.T) {
	defer ExpectPanic("Cannot use the provided dependency 'book_creator' of type 'mocks.BookCreator' as 'mocks.BookStorage' in the constructor function call", t)

	cont := CreateContainer()
	cont.AddNewMethod("anotherFinder", mocks.NewBookFinder, "book_creator", "book_storage")

	var finder mocks.BookFinder
	cont.Scan("anotherFinder", &finder)
}

func TestWrongDependencyRequested(t *testing.T) {
	defer ExpectPanic("Unknown dependency 'lala'", t)
	cont := CreateContainer()

	cont.Scan("lala", nil)
}

func TestIncompatibleInterfaces(t *testing.T) {
	defer ExpectPanic("Cannot use the provided dependency 'incompatible_cache' of type 'mocks.IncompatibleCache' as 'mocks.Cache' in the constructor function call", t)
	cont := CreateContainer()
	cont.AddNewMethod("incompatible_cache", mocks.NewIncompatibleCache)
	cont.AddNewMethod(
		"wrong_downloader",
		mocks.NewBookDownloader,
		"incompatible_cache",
		"book_link_provider",
		"book_finder",
		"web_fetcher",
	)

	cont.Scan("incompatible_cache", nil)
}

func TestCheckFailingForWrongLazyDependencies(t *testing.T) {
	defer ExpectPanic("Cannot convert created value of type 'int' to expected destination value 'BookCreator' for createdDependency declaration wrong_book_creator", t)
	cont := CreateContainer()
	cont.AddConstructor("wrong_book_finder", func(c Container) (interface{}, error) {
		var bc mocks.BookCreator
		c.Scan("wrong_book_creator", &bc)

		var bs mocks.BookStorage
		c.Scan("book_storage", &bs)

		return mocks.NewBookFinder(bs, bc), nil
	})

	cont.Check()
}

func TestCorrectCheck(t *testing.T) {
	cont := CreateContainer()
	cont.Check()
}
