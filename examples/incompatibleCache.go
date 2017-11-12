package examples

//IncompatibleCache simulates a mock that doesn't implement the Cache interface
type IncompatibleCache struct{}

//Simulates constructor
func NewIncompatibleCache() IncompatibleCache {
	return IncompatibleCache{}
}

//Not compatible func declaration
func (dependencyCache IncompatibleCache) Cache(book Book) bool {
	return false
}
