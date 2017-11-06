package examples

type IncompatibleCache struct {}

func NewIncompatibleCache() IncompatibleCache {
	return IncompatibleCache{}
}

func (cacheService IncompatibleCache) Cache(book Book) bool {
	return false
}
