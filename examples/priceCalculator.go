package examples

type PriceCalculator struct {
	priceFinder     PriceFinder
	priceMultiplier func(inputPrice int) int
	priceDiscounter func(inputPrice int) int
}

func NewPriceCalculator(priceFinder PriceFinder, priceMultiplier func(inputPrice int) int) *PriceCalculator {
	return &PriceCalculator{
		priceFinder:     priceFinder,
		priceMultiplier: priceMultiplier,
		priceDiscounter: func(inputPrice int) int {
			return inputPrice
		},
	}
}

func NewPriceDoubler() func(inputPrice int) int {
	return func(inputPrice int) int {
		return inputPrice * 2
	}
}

func NewPriceTripleMultiplier() func(inputPrice int) int {
	return func(inputPrice int) int {
		return inputPrice * 3
	}
}

func (pc *PriceCalculator) CalculateBookPrice(bookId string) int {
	price := pc.priceFinder(bookId)
	price = pc.priceMultiplier(price)
	price = pc.priceDiscounter(price)

	return price
}

func (pc *PriceCalculator) SetDiscounter(discounter func(inputPrice int) int) {
	pc.priceDiscounter = discounter
}
