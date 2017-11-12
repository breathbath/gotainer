package examples

//PriceCalculator simulates a price computation engine
type PriceCalculator struct {
	priceFinder     PriceFinder
	priceMultiplier func(inputPrice int) int
	priceDiscounter func(inputPrice int) int
}

//Constructor
func NewPriceCalculator(priceFinder PriceFinder, priceMultiplier func(inputPrice int) int) *PriceCalculator {
	return &PriceCalculator{
		priceFinder:     priceFinder,
		priceMultiplier: priceMultiplier,
		priceDiscounter: func(inputPrice int) int {
			return inputPrice
		},
	}
}

//NewPriceDoubler multiplies a price by 2
func NewPriceDoubler() func(inputPrice int) int {
	return func(inputPrice int) int {
		return inputPrice * 2
	}
}

//NewPriceTripleMultiplier multiplies a price by 3
func NewPriceTripleMultiplier() func(inputPrice int) int {
	return func(inputPrice int) int {
		return inputPrice * 3
	}
}

//CalculateBookPrice applies all prices modifies and returns result
func (pc *PriceCalculator) CalculateBookPrice(bookId string) int {
	price := pc.priceFinder(bookId)
	price = pc.priceMultiplier(price)
	price = pc.priceDiscounter(price)

	return price
}

//SetDiscounter simulates an optional price modifier
func (pc *PriceCalculator) SetDiscounter(discounter func(inputPrice int) int) {
	pc.priceDiscounter = discounter
}
