package container

type Container interface {
	AddConstructor(id string, constructor Constructor)
	AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string)
	Scan(id string, dest interface{})
	ScanNonCached(id string, dest interface{})
	Get(id string, isCached bool) interface{}
	Check()
}
