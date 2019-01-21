package container

import (
	"errors"
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestWrongMixedDependenciesInStaticCall(t *testing.T) {
	defer ExpectPanic(
		t,
		"Cannot use the provided dependency 'book_creator' of type 'mocks.BookCreator' as 'mocks.BookStorage' in the Constr function call [check 'anotherFinder' service];\n"+
			"Cannot use the provided dependency 'book_storage' of type 'mocks.BookStorage' as 'mocks.BookCreator' in the Constr function call [check 'anotherFinder' service]",
	)

	cont := CreateContainer()
	cont.AddNewMethod("anotherFinder", mocks.NewBookFinder, "book_creator", "book_storage")

	var finder mocks.BookFinder
	cont.Scan("anotherFinder", &finder)
}

func TestWrongDependencyRequested(t *testing.T) {
	defer ExpectPanic(t, "Unknown dependency 'lala'")
	cont := CreateContainer()

	cont.Scan("lala", nil)
}

func TestWrongDependencyGetFromSecureMethod(t *testing.T) {
	cont := CreateContainer()

	_, err := cont.GetSecure("lala", true)
	if err.Error() != "Unknown dependency 'lala'" {
		t.Errorf("Unexpected error: %v, expected error was Unknown dependency 'lala'", err)
	}
}

func TestWrongDependencyScannedFromSecureMethod(t *testing.T) {
	cont := CreateContainer()
	err := cont.ScanSecure("lala", true, nil)
	if err.Error() != "Unknown dependency 'lala'" {
		t.Errorf("Unexpected error: %v, expected error was Unknown dependency 'lala'", err)
	}
}

func TestIncompatibleInterfaces(t *testing.T) {
	defer ExpectPanic(t, "Cannot use the provided dependency 'incompatible_cache' of type 'mocks.IncompatibleCache' as 'mocks.Cache' in the Constr function call [check 'wrong_downloader' service]")
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

	var bd mocks.BookDownloader

	cont.Scan("wrong_downloader", &bd)
}

func TestCheckFailingForWrongLazyDependencies(t *testing.T) {
	defer ExpectPanic(t, "Cannot convert created value of type 'int' to expected destination value 'BookCreator' for createdDependency declaration wrong_book_creator [check 'wrong_book_creator' service]")
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

func TestCheckFailingForInvalidGarbageCollectionDeclaration(t *testing.T) {
	cont := CreateContainer()
	garbageCollector := func(service interface{}) error {
		return nil
	}
	cont.AddGarbageCollectFunc("some_unknown_service", garbageCollector)

	err := cont.CollectGarbage()
	expectedErrorText := "Garbage collection errors: Unknown dependency 'some_unknown_service'"
	if err.Error() != expectedErrorText {
		t.Errorf("Unexpected error: %v, expected error was %s", err, expectedErrorText)
	}
}

func TestFailureOnCustomConstructorError(t *testing.T) {
	cont := CreateContainer()
	cont.AddConstructor("some_failing_constructor", func(c Container) (interface{}, error) {
		return nil, errors.New("Something bad has happened")
	})
	defer ExpectPanic(t, "Something bad has happened [check 'some_failing_constructor' service]")
	cont.Get("some_failing_constructor", true)
}

func TestWrongArgumentsCountInNewMethod(t *testing.T) {
	cont := CreateContainer()
	defer ExpectPanic(t, "The function requires 2 arguments, but 3 arguments are provided [check 'wrong_arg_count_dependency' service]")
	cont.AddNewMethod("wrong_arg_count_dependency", mocks.NewBookFinder, "book_storage", "book_creator", "Config")
}

func TestNewMethodIsNotFunc(t *testing.T) {
	cont := CreateContainer()
	defer ExpectPanic(t, "A function is expected rather than 'string' [check 'wrong_new_method_dependency' service]")
	cont.AddNewMethod("wrong_new_method_dependency", "non_func")
}

func TestValidationOfReturnsCountInNewMethod(t *testing.T) {
	cont := CreateContainer()
	someBadNewFunc := func() (int, error, bool) {
		return 1, errors.New("some error"), true
	}
	defer ExpectPanic(t, "Constr function should return 1 or 2 values, but 3 values are returned [check 'wrong_return_count_dependency' service]")
	cont.AddNewMethod("wrong_return_count_dependency", someBadNewFunc)
}

func TestNewMethodWithTwoReturnsHasAtLeastOneError(t *testing.T) {
	cont := CreateContainer()
	someBadNewFunc := func() (int, bool) {
		return 1, true
	}
	defer ExpectPanic(t, "Constr function with 2 returned values should return at least one error interface [check 'wrong_two_returns_with_no_error_dependency' service]")
	cont.AddNewMethod("wrong_two_returns_with_no_error_dependency", someBadNewFunc)
}

func TestNewMethodReturnError(t *testing.T) {
	cont := CreateContainer()
	newMethodWithError := func() (error, interface{}) {
		return errors.New("New method failed for some reason"), nil
	}
	cont.AddNewMethod("some_failing_newMethod", newMethodWithError)
	defer ExpectPanic(t, "New method failed for some reason [check 'some_failing_newMethod' service]")
	cont.Get("some_failing_newMethod", true)
}

func TestNonPointerVariableDestination(t *testing.T) {
	cont := CreateContainer()
	var url string
	defer ExpectPanic(t, "Please provide a pointer variable rather than a value [check 'static_files_url' service]")
	cont.Scan("static_files_url", url)
}

func TestNonInitialisedPointerVariableDestination(t *testing.T) {
	cont := CreateContainer()
	var url *string
	defer ExpectPanic(t, "Please provide an initialized variable rather than a non-initialised pointer variable [check 'static_files_url' service]")
	cont.Scan("static_files_url", url)
}
