package container

type Container interface {
	AddNoArgumentsConstructor(id string, constructor NoArgumentsConstructor)
	AddConstructor(id string, constructor ArgumentsConstructor)
	AddTypedConstructor(id string, typedConstructor interface{}, constructorArgumentNames ...string)
	GetTypedService(id string, dest interface{})
	GetService(id string) interface{}
}
