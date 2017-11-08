package examples

type PriceFinder func(bookId string) int


func NewBooksPriceFinder(bookPrices map[string] int, books []Book) PriceFinder {
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

func GetBookPrices() map[string] int {
	return map[string] int {"1" : 100, "2": 200}
}

func GetAllBooks() []Book {
	return []Book {Book{"1", "Book1", "Author1"}, Book{"2", "Book2", "Author2"}}
}