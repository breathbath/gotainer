package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
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
