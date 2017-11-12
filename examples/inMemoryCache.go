package examples

//Cache for your books
type InMemoryCache struct {
	cache map[string]Book
}

//Main constructor
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string]Book)}
}

//Caches a book by its id
func (dependencyCache *InMemoryCache) Cache(book Book) {
	dependencyCache.cache[book.Id] = book
}
