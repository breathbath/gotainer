package mocks

//InMemoryCache cache for your books
type InMemoryCache struct {
	cache map[string]Book
}

//NewInMemoryCache InMemoryCache constructor
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string]Book)}
}

//Cache saves a book by id in a cache
func (dependencyCache *InMemoryCache) Cache(book Book) {
	dependencyCache.cache[book.Id] = book
}
