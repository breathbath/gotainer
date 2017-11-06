package examples

type InMemoryCache struct {
	cache map[string] Book
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string]Book)}
}

func (cacheService *InMemoryCache) Cache(book Book) {
	cacheService.cache[book.Id] = book
}