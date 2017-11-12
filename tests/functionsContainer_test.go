package tests

import (
	"github.com/breathbath/gotainer/container"
	"github.com/breathbath/gotainer/examples"
	"testing"
)

func TestFunctionalDependency(t *testing.T) {
	cont := PrepareContainer()
	var priceCalculator examples.PriceCalculator
	cont.Scan("price_calculator_double", &priceCalculator)

	AssertPrice(200, "1", priceCalculator, t)

	cont.Scan("price_calculator_triple", &priceCalculator)
	AssertPrice(600, "2", priceCalculator, t)
}

func AssertPrice(expectedPrice int, bookId string, priceCalculator examples.PriceCalculator, t *testing.T) {
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

	cont.AddConstructor("price_calculator_double_discount", func(c container.Container) (interface{}, error) {
		var priceCalculatorDouble examples.PriceCalculator
		c.Scan("price_calculator_double", &priceCalculatorDouble)
		priceCalculatorDouble.SetDiscounter(discountFunc)

		return priceCalculatorDouble, nil
	})

	var priceCalculator examples.PriceCalculator
	cont.Scan("price_calculator_double_discount", &priceCalculator)

	AssertPrice(180, "1", priceCalculator, t)
}

func PrepareContainer() container.Container {
	cont := examples.CreateContainer()
	cont.AddNewMethod("book_prices", examples.GetBookPrices)
	cont.AddNewMethod("books", examples.GetAllBooks)
	cont.AddNewMethod("price_finder", examples.NewBooksPriceFinder, "book_prices", "books")

	cont.AddNewMethod("price_doubler", examples.NewPriceDoubler)
	cont.AddNewMethod(
		"price_calculator_double",
		examples.NewPriceCalculator,
		"price_finder",
		"price_doubler",
	)

	cont.AddNewMethod("price_tripleMultiplier", examples.NewPriceTripleMultiplier)
	cont.AddNewMethod(
		"price_calculator_triple",
		examples.NewPriceCalculator,
		"price_finder",
		"price_tripleMultiplier",
	)

	return cont
}
