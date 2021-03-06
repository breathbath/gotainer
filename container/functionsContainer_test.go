package container

import (
	"fmt"
	"github.com/breathbath/gotainer/container/mocks"
	"strings"
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

func TestDuplicatedGarbageCollectFuncs(t *testing.T) {
	cont := PrepareContainer()

	funcCalledNumber := 0
	garbageCollector1 := func(service interface{}) error {
		funcCalledNumber = 1
		return nil
	}
	garbageCollector2 := func(service interface{}) error {
		funcCalledNumber = 2
		return nil
	}

	cont.AddGarbageCollectFunc("in_memory_cache", garbageCollector1)
	cont.AddGarbageCollectFunc("in_memory_cache", garbageCollector2)

	err := cont.CollectGarbage()
	if err != nil {
		t.Errorf("Unexpected error %v during the garbage collection", err)
		return
	}

	if funcCalledNumber != 1 {
		t.Errorf("First garbage collection func should be called rather than second which has a conflicting name")
	}
}

func TestGarbageCollectionLoop(t *testing.T) {
	funcsCollection := NewGarbageCollectorFuncs()
	funcCalls := []string{}
	funcsCollection.Add("func1", func(service interface{}) error {
		funcCalls = append(funcCalls, "func1")
		return nil
	})
	funcsCollection.Add("func2", func(service interface{}) error {
		funcCalls = append(funcCalls, "func2")
		return nil
	})
	funcsCollection.Add("func3", func(service interface{}) error {
		funcCalls = append(funcCalls, "func3")
		return nil
	})

	funcsCollection.Range(func(gcName string, f GarbageCollectorFunc) bool {
		f(nil)

		if gcName == "func2" {
			return false
		}

		return true
	})

	actualCalledFuncs := strings.Join(funcCalls, ",")
	if strings.Join(funcCalls, ",") != "func1,func2" {
		t.Errorf("It was expected that only func1 and func2 are called but %s functions were called", actualCalledFuncs)
	}
}

func TestGarbageCollectionCallOrder(t *testing.T) {
	garbageCollectionServices := []string{
		"book_prices",
		"books",
		"price_finder",
		"price_doubler",
		"price_calculator_double",
		"price_tripleMultiplier",
	}

	actualCalls := []string{}

	cont := PrepareContainer()

	gcFuncBuilder := func(gcServiceName string) GarbageCollectorFunc {
		return func(service interface{}) error {
			actualCalls = append(actualCalls, gcServiceName)
			return nil
		}
	}
	for _, gcServiceName := range garbageCollectionServices {
		gcFunc := gcFuncBuilder(gcServiceName)
		cont.AddGarbageCollectFunc(gcServiceName, gcFunc)
	}

	err := cont.CollectGarbage()
	if err != nil {
		t.Error(err)
		return
	}

	expectedCallsChain := strings.Join(garbageCollectionServices, ",")
	actualCallsChain := strings.Join(actualCalls, ",")

	if expectedCallsChain != actualCallsChain {
		t.Errorf(`
-Expected garbage collection calls 
+Actual garbage collection calls
- [%s]
+ [%s]
are not equal`,
			expectedCallsChain,
			actualCallsChain)
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

	expectedError := "Garbage collection errors: Error 1, Error 2"
	if err.Error() == expectedError {
		return
	}

	t.Errorf("Garbage collect function should return '%s' but '%v' is returned", expectedError, err)
}

func TestGarbageCollectionForUnknownService(t *testing.T) {
	cont := PrepareContainer()
	garbageCollector := func(service interface{}) error {
		return nil
	}

	cont.AddGarbageCollectFunc("some_unknown_service", garbageCollector)
	err := cont.CollectGarbage()
	expectedError := "Garbage collection errors: Unknown dependency 'some_unknown_service'"
	if err.Error() != expectedError {
		t.Errorf("Unexpected error: '%v' expected error was '%s'", err, expectedError)
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
