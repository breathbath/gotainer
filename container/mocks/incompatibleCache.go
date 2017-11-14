package mocks

//IncompatibleCache simulates a mock that doesn't implement the Cache interface
type IncompatibleCache struct{}

//NewIncompatibleCache IncompatibleCache constructor
func NewIncompatibleCache() IncompatibleCache {
	return IncompatibleCache{}
}

//Cache a incompatible func declaration
func (dependencyCache IncompatibleCache) Cache(book Book) bool {
	return false
}
