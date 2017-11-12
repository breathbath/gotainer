package examples

//Caching layer for book entities
type Cache interface {
	Cache(book Book)
}
