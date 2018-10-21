package container

import (
	"fmt"
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestFunctionalDependency(t *testing.T) {
	cont := PrepareContainer()
	var priceCalculator mocks.PriceCalculator
	cont.Scan("price_calculator_double", &priceCalculator)

	AssertPrice(200, "1", priceCalculator, t)

	cont.Scan("price_calculator_triple", &priceCalculator)
	AssertPrice(600, "2", priceCalculator, t)
}

func AssertPrice(expectedPrice int, bookId string, priceCalculator mocks.PriceCalculator, t *testing.T) {
	receivedPrice := priceCalculator.CalculateBookPrice(bookId)
	if receivedPrice != expectedPrice {
		t.Errorf(
			"Price calculated incorrectly, expected price '%d', provided price '%d', book id '%s'",
			expectedPrice,
			receivedPrice,
			bookId,
		)
	}
}

func TestFunctionalDependencyWithSetter(t *testing.T) {
	cont := PrepareContainer()
	discountFunc := func(inputPrice int) int {
		return inputPrice - inputPrice/10
	}

	cont.AddConstructor("price_calculator_double_discount", func(c Container) (interface{}, error) {
		var priceCalculatorDouble mocks.PriceCalculator
		c.Scan("price_calculator_double", &priceCalculatorDouble)
		priceCalculatorDouble.SetDiscounter(discountFunc)

		return priceCalculatorDouble, nil
	})

	var priceCalculator mocks.PriceCalculator
	cont.Scan("price_calculator_double_discount", &priceCalculator)

	AssertPrice(180, "1", priceCalculator, t)
}

func TestGarbageCollectionSuccess(t *testing.T) {
	cont := PrepareContainer()
	wasCalled := false
	garbageCollector := func(service interface{}) error {
		priceDoubler := service.(func(inputPrice int) int)
		priceDoubler(1)
		wasCalled = true

		return nil
	}

	cont.AddGarbageCollectFunc("price_doubler", garbageCollector)

	if wasCalled {
		t.Error("Garbage collect function should not have been called upon initialisation")
		return
	}

	err := cont.CollectGarbage()
	if err != nil {
		t.Errorf("Unexpected error %v during the garbage collection", err)
		return
	}

	if !wasCalled {
		t.Error("The registered garbage collect function should have been called")
	}
}

func TestGarbageCollectionFailures(t *testing.T) {
	cont := PrepareContainer()

	garbageCollector1 := func(service interface{}) error {
		return fmt.Errorf("Error 1")
	}
	cont.AddGarbageCollectFunc("book_prices", garbageCollector1)

	garbageCollector2 := func(service interface{}) error {
		return fmt.Errorf("Error 2")
	}
	cont.AddGarbageCollectFunc("books", garbageCollector2)

	garbageCollector3 := func(service interface{}) error {
		return nil
	}
	cont.AddGarbageCollectFunc("price_finder", garbageCollector3)

	err := cont.CollectGarbage()

	expectedError1 := "Garbage collection errors: Error 1, Error 2"
	expectedError2 := "Garbage collection errors: Error 2, Error 1"
	if err.Error() == expectedError1 || err.Error() == expectedError2 {
		return
	}

	t.Errorf("Garbage collect function should return '%s' or '%s' but '%v' is returned", expectedError1, expectedError2, err)
}

func TestGarbageCollectionForUnknownService(t *testing.T) {
	defer ExpectPanic("Unknown dependency 'some_unknown_service'", t)
	cont := PrepareContainer()
	garbageCollector := func(service interface{}) error {
		return nil
	}

	cont.AddGarbageCollectFunc("some_unknown_service", garbageCollector)
	err := cont.CollectGarbage()
	if err != nil {
		t.Errorf("Unexpected error %v during the garbage collection", err)
	}
}

func PrepareContainer() Container {
	cont := CreateContainer()
	cont.AddNewMethod("book_prices", mocks.GetBookPrices)
	cont.AddNewMethod("books", mocks.GetAllBooks)
	cont.AddNewMethod("price_finder", mocks.NewBooksPriceFinder, "book_prices", "books")

	cont.AddNewMethod("price_doubler", mocks.NewPriceDoubler)
	cont.AddNewMethod(
		"price_calculator_double",
		mocks.NewPriceCalculator,
		"price_finder",
		"price_doubler",
	)

	cont.AddNewMethod("price_tripleMultiplier", mocks.NewPriceTripleMultiplier)
	cont.AddNewMethod(
		"price_calculator_triple",
		mocks.NewPriceCalculator,
		"price_finder",
		"price_tripleMultiplier",
	)

	return cont
}
