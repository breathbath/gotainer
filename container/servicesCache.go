package container

type servicesCache map[string]interface{}

func newServicesCache() servicesCache {
	return make(map[string]interface{})
}

func (sC servicesCache) Get(id string) (interface{}, bool){
	service, ok := sC[id]
	return service, ok
}
