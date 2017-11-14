package mocks

//PriceFinder is intended to find prices of books
type PriceFinder func(bookId string) int

//NewBooksPriceFinder PriceFinder constructor
func NewBooksPriceFinder(bookPrices map[string]int, books []Book) PriceFinder {
	return func(bookId string) int {
		price, ok := bookPrices[bookId]
		if !ok {
			return 0
		}

		for _, book := range books {
			if book.Id == bookId {
				return price
			}
		}
		return 0
	}
}

//GetBookPrices gives a list of prices for books
func GetBookPrices() map[string]int {
	return map[string]int{"1": 100, "2": 200}
}
