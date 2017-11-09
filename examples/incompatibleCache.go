package examples

type IncompatibleCache struct {}

func NewIncompatibleCache() IncompatibleCache {
	return IncompatibleCache{}
}

func (dependencyCache IncompatibleCache) Cache(book Book) bool {
	return false
}
