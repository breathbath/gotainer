package container

//dependencyCache is a in-memory cache for all declared services
type dependencyCache map[string]interface{}

func newDependencyCache() dependencyCache {
	return make(map[string]interface{})
}

func (sC dependencyCache) Get(id string) (interface{}, bool) {
	dep, ok := sC[id]
	return dep, ok
}

func (sC dependencyCache) Set(id string, dep interface{}) bool {
	sC[id] = dep
	return true
}
