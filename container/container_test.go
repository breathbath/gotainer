package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"os"
	"testing"
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
	expectedURL := "http://static.me/FirstBook"
	if url != expectedURL {
		t.Errorf(
			"Unexpected book url '%s' fetched, expected url is '%s'.",
			url,
			expectedURL,
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

func TestCheckNotFails(t *testing.T) {
	cont := CreateContainer()
	cont.Check()
}

func TestExistsFunction(t *testing.T) {
	cont := CreateContainer()
	if !cont.Exists("book_storage") {
		t.Errorf("The service '%s' should exist", "book_storage")
	}

	if cont.Exists("some_non_existing_service") {
		t.Errorf("The service '%s' should not exist", "some_non_existing_service")
	}
}

func TestSettingConstructor(t *testing.T) {
	cont := CreateContainer()
	cont.SetConstructor("someService", func(c Container) (i interface{}, e error) {
		return "currentName", nil
	})

	cont.SetConstructor("someService", func(c Container) (i interface{}, e error) {
		return "overriddenName", nil
	})

	actualValue := cont.Get("someService", true).(string)
	if actualValue != "overriddenName" {
		t.Errorf(
			"Expected service 'someService' is not equal to returned '%s' value. Second constructor declaration for 'someService' should override an existing declaration",
			actualValue,
		)
	}

	cont.AddConstructor("addedService", func(c Container) (i interface{}, e error) {
		return "addedServiceName", nil
	})
	cont.SetConstructor("addedService", func(c Container) (i interface{}, e error) {
		return "overriddenAddedServiceName", nil
	})

	actualValue = cont.Get("addedService", true).(string)
	if actualValue != "overriddenAddedServiceName" {
		t.Errorf(
			"Expected service 'overriddenAddedServiceName' is not equal to returned '%s' value. Second constructor declaration for 'addedService' should override an existing declaration",
			actualValue,
		)
	}
}

func TestSettingNewMethod(t *testing.T) {
	cont := CreateContainer()
	newMethod1 := func() int {
		return 1
	}
	newMethod2 := func() int {
		return 2
	}

	cont.SetNewMethod("counter", newMethod1)
	cont.SetNewMethod("counter", newMethod2)

	actualValue := cont.Get("counter", true).(int)
	if actualValue != 2 {
		t.Errorf(
			"Expected service '2' is not equal to returned '%d' value. Second new method declaration for 'counter' should override an existing declaration",
			actualValue,
		)
	}
}

/**
In my prod code I would do
	cont := CreateContainerForPaymentsCase()
	registrator := cont.Get("registrator", true).(Registrator)
	err := registrator.RegisterUser("client_1")
 */
//now I can test how registrator handles error responses from the payment gateway
func TestRegistrationPaymentGatewayFailure(t *testing.T) {
	cont := CreateContainerForPaymentsCase()
	//here I am replacing the NewRealPaymentGateway with the NewFailingPaymentGateway as both return
	//implementation of PaymentGateway interface so no failure will happen
	cont.SetNewMethod("paymentGateway", mocks.NewFailingPaymentGateway)

	//now the registrator contains instead of RealPaymentGateway the FailingPaymentGateway so I can test my case
	registrator := cont.Get("registrator", true).(mocks.Registrator)
	err := registrator.RegisterUser("client_1")

	//I expect that the registrator just forwards the original api error and of course that the replacement took place
	if err.Error() != "Cannot connect to external api" {
		t.Errorf(
			"Unexpected error is returned: %v, expected exception is 'Cannot connect to external api'",
			err,
		)
	}
}

func CreateContainerForPaymentsCase() *RuntimeContainer {
	cont := NewRuntimeContainer()
	cont.AddConstructor("secretKey", func(c Container) (i interface{}, e error) {
		//even if we replace env var with some fake key in testing env
		// we still will do real call to the external payment gateway
		return os.Getenv("PAYMENT_SECRET_KEY"), nil
	})
	cont.AddNewMethod("paymentGateway", mocks.NewRealPaymentGateway, "secretKey")
	cont.AddNewMethod("registrator", mocks.NewRegistrator, "paymentGateway")

	return cont
}
