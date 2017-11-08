package tests

import (
	"testing"
	"github.com/breathbath/gotainer/examples"
	"github.com/breathbath/gotainer/container"
	"fmt"
)

func TestFunctionalDependency(t *testing.T) {
	cont := PrepareContainer()
	var priceCalculator examples.PriceCalculator
	cont.GetTypedService("price_calculator_double", &priceCalculator)

	AssertPrice(200, "1", priceCalculator, t)

	cont.GetTypedService("price_calculator_triple", &priceCalculator)
	AssertPrice(600, "2", priceCalculator, t)
}

func AssertPrice(expectedPrice int, bookId string, priceCalculator examples.PriceCalculator, t *testing.T) {
	receivedPrice := priceCalculator.CalculateBookPrice(bookId)
	if receivedPrice != expectedPrice {
		fmt.Errorf(
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

	cont.AddConstructor("price_calculator_double_discount", func(c container.Container) (interface{}, error) {
		var priceCalculatorDouble examples.PriceCalculator
		c.GetTypedService("price_calculator_double", &priceCalculatorDouble)
		priceCalculatorDouble.SetDiscounter(discountFunc)

		return priceCalculatorDouble, nil
	})

	var priceCalculator examples.PriceCalculator
	cont.GetTypedService("price_calculator_double_discount", &priceCalculator)

	AssertPrice(180, "1", priceCalculator, t)
}

func PrepareContainer() container.Container {
	cont := examples.CreateContainer()
	cont.AddTypedConstructor("book_prices", examples.GetBookPrices)
	cont.AddTypedConstructor("books", examples.GetAllBooks)
	cont.AddTypedConstructor("price_finder", examples.NewBooksPriceFinder, "book_prices", "books")

	cont.AddTypedConstructor("price_doubler", examples.NewPriceDoubler)
	cont.AddTypedConstructor(
		"price_calculator_double",
		examples.NewPriceCalculator,
		"price_finder",
		"price_doubler",
	)

	cont.AddTypedConstructor("price_tripleMultiplier", examples.NewPriceTripleMultiplier)
	cont.AddTypedConstructor(
		"price_calculator_triple",
		examples.NewPriceCalculator,
		"price_finder",
		"price_tripleMultiplier",
	)

	return cont
}
