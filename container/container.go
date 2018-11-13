package container

//Container main interface for registering and fetching Services
type Container interface {
	AddConstructor(id string, constructor Constructor)
	AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string)
	Scan(id string, dest interface{})
	ScanNonCached(id string, dest interface{})
	ScanSecure(id string, isCached bool, dest interface{}) error
	Get(id string, isCached bool) interface{}
	GetSecure(id string, isCached bool) (interface{}, error)
	Check()
	Exists(id string) bool
	AddGarbageCollectFunc(serviceName string, gcFunc GarbageCollectorFunc)
	CollectGarbage() error
}

//MergeableContainer containers that support merging
type MergeableContainer interface {
	Merge(c MergeableContainer)
	getConstructors() map[string]Constructor
	getNewFuncConstructors() map[string]NewFuncConstructor
	getCache() dependencyCache
	getEventsContainer() EventsContainer
}
