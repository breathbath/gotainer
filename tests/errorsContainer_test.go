package tests

import (
	"testing"
	"github.com/breathbath/gotainer/examples"
	"github.com/breathbath/gotainer/container"
)

func TestWrongMixedDependenciesInStaticCall(t *testing.T) {
	defer ExpectPanic("Cannot use the provided dependency 'book_creator' of type 'examples.BookCreator' as 'examples.BookStorage' in the constructor function call", t)

	cont := examples.CreateContainer()
	cont.AddNewMethod("anotherFinder", examples.NewBookFinder, "book_creator", "book_storage")

	var finder examples.BookFinder
	cont.Scan("anotherFinder", &finder)
}

func TestWrongDependencyRequested(t *testing.T) {
	defer ExpectPanic("Unknown dependency 'lala'", t)
	cont := examples.CreateContainer()

	cont.Scan("lala", nil)
}

func TestIncompatibleInterfaces(t *testing.T) {
	defer ExpectPanic("Cannot use the provided dependency 'incompatible_cache' of type 'examples.IncompatibleCache' as 'examples.Cache' in the constructor function call", t)
	cont := examples.CreateContainer()
	cont.AddNewMethod("incompatible_cache", examples.NewIncompatibleCache)
	cont.AddNewMethod(
		"wrong_downloader",
		examples.NewBookDownloader,
		"incompatible_cache",
		"book_link_provider",
		"book_finder",
		"web_fetcher",
	)

	cont.Scan("incompatible_cache", nil)
}

func TestCheckFailingForWrongLazyDependencies(t *testing.T) {
	defer ExpectPanic("Cannot convert created value of type 'int' to expected destination value 'BookCreator' for createdDependency declaration wrong_book_creator", t)
	cont := examples.CreateContainer()
	cont.AddConstructor("wrong_book_finder", func(c container.Container) (interface{}, error) {
		var bc examples.BookCreator
		c.Scan("wrong_book_creator", &bc)

		var bs examples.BookStorage
		c.Scan("book_storage", &bs)

		return examples.NewBookFinder(bs, bc), nil
	})

	cont.Check()
}

func TestCorrectCheck(t *testing.T) {
	cont := examples.CreateContainer()
	cont.Check()
}