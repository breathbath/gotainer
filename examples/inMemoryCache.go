package examples

type InMemoryCache struct {
	cache map[string]Book
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string]Book)}
}

func (dependencyCache *InMemoryCache) Cache(book Book) {
	dependencyCache.cache[book.Id] = book
}
