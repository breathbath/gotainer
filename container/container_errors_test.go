package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
	"errors"
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

func TestFailureOnCustomConstructorError(t *testing.T) {
	cont := CreateContainer()
	cont.AddConstructor("some_failing_constructor", func(c Container) (interface{}, error) {
		return nil, errors.New("Something bad has happened")
	})
	defer ExpectPanic("Something bad has happened", t)
	cont.Get("some_failing_constructor", true)
}

func TestWrongArgumentsCountInNewMethod(t *testing.T) {
	cont := CreateContainer()
	defer ExpectPanic("The function requires 2 dependencies, but 3 arguments are provided", t)
	cont.AddNewMethod("wrong_arg_count_dependency", mocks.NewBookFinder, "book_storage", "book_creator", "config")
}

func TestNewMethodIsNotFunc(t *testing.T) {
	cont := CreateContainer()
	defer ExpectPanic("Destination object should be a constructor function rather than string", t)
	cont.AddNewMethod("wrong_arg_count_dependency", "non_func")
}

func TestValidationOfReturnsCountInNewMethod(t *testing.T) {
	cont := CreateContainer()
	someBadNewFunc := func() (int, error, bool){
		return 1, errors.New("some error"), true
	}
	defer ExpectPanic("constructor function should return 1 or 2 values, but 3 values are returned", t)
	cont.AddNewMethod("wrong_return_count_dependency", someBadNewFunc)
}

func TestNewMethodWithTwoReturnsHasAtLeastOneError(t *testing.T) {
	cont := CreateContainer()
	someBadNewFunc := func() (int, bool){
		return 1, true
	}
	defer ExpectPanic("constructor function with 2 returned values should return at least one error interface", t)
	cont.AddNewMethod("wrong_two_returns_with_no_error_dependency", someBadNewFunc)
}

func TestNewMethodReturnError(t *testing.T) {
	cont := CreateContainer()
	newMethodWithError := func() (error, interface{}) {
		return errors.New("New method failed for some reason"), nil
	}
	cont.AddNewMethod("some_failing_newMethod", newMethodWithError)
	defer ExpectPanic("New method failed for some reason", t)
	cont.Get("some_failing_newMethod", true)
}