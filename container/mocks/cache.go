package mocks

//Cache caching layer for book entities
type Cache interface {
	Cache(book Book)
}
