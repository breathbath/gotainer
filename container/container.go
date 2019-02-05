package container

//Container main interface for registering and fetching Services
type Container interface {
	AddConstructor(id string, constructor Constructor) error
	AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) error
	Scan(id string, dest interface{})
	ScanNonCached(id string, dest interface{})
	ScanSecure(id string, isCached bool, dest interface{}) error
	Get(id string, isCached bool) interface{}
	GetSecure(id string, isCached bool) (interface{}, error)
	Check() error
	Exists(id string) bool
	AddGarbageCollectFunc(serviceName string, gcFunc GarbageCollectorFunc)
	CollectGarbage() error
	SetConstructor(id string, constructor Constructor)
	SetNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) error
}

//MergeableContainer containers that support merging
type MergeableContainer interface {
	Merge(c MergeableContainer) error
	getConstructors() map[string]Constructor
	getNewFuncConstructors() map[string]NewFuncConstructor
	getCache() dependencyCache
	getEventsContainer() EventsContainer
}
