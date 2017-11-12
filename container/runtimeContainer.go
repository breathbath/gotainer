package container

import (
	"errors"
	"fmt"
)

//RuntimeContainer creates services at runtime with registered callbacks
type RuntimeContainer struct {
	constructors    map[string]Constructor
	cache           dependencyCache
	eventsContainer *EventsContainer
}

//NewRuntimeContainer creates container
func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{constructors: make(map[string]Constructor), cache: newDependencyCache(), eventsContainer: NewEventsContainer()}
}

//AddConstructor registers a callback to create a service identified by id
func (rc *RuntimeContainer) AddConstructor(id string, constructor Constructor) {
	rc.constructors[id] = constructor
}

//AddNewMethod converts a New service method to a valid callback constructor
func (rc *RuntimeContainer) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	rc.constructors[id] = convertNewMethodToConstructor(rc, typedConstructor, constructorArgumentNames)
}

//AddDependencyObserver registers service that will receive dependencies it is interested in
func (rc *RuntimeContainer) AddDependencyObserver(eventName, observerId string, observer interface{}) {
	rc.eventsContainer.addDependencyObserver(eventName, observerId, observer)
}

//RegisterDependencyEvent notifies observers about added dependencies
func (rc *RuntimeContainer) RegisterDependencyEvent(eventName, dependencyName string) {
	rc.eventsContainer.registerDependencyEvent(eventName, dependencyName)
}

//Scan copies a service identified by id into a typed destination (its a pointer reference)
func (rc *RuntimeContainer) Scan(id string, dest interface{}) {
	baseValue := rc.Get(id, true)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

//ScanNonCached creates a service every time rc method is called
func (rc *RuntimeContainer) ScanNonCached(id string, dest interface{}) {
	baseValue := rc.Get(id, false)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

//Get fetches a service in a return argument
func (rc *RuntimeContainer) Get(id string, isCached bool) interface{} {
	dependency, ok := rc.cache.Get(id)
	if ok && isCached {
		return dependency
	}

	constructorFunc, ok := rc.constructors[id]
	if !ok {
		errStr := fmt.Sprintf("Unknown dependency '%s'", id)
		panic(errors.New(errStr))
	}

	result, err := constructorFunc(rc)
	if err != nil {
		panic(err)
	}

	rc.eventsContainer.notifyObserverAboutDependency(*rc, id, result)

	rc.cache.Set(id, result)

	return result
}

//Check ensures that all runtime dependencies are created correctly
func (rc *RuntimeContainer) Check() {
	for dependencyName := range rc.constructors {
		rc.Get(dependencyName, false)
	}
}

//Merge allows to merge containers
func (rc *RuntimeContainer) Merge(c MergeableContainer) {
	for keyConstructor, constr := range c.getConstructors() {
		rc.constructors[keyConstructor] = constr
	}

	for keyCache, cache := range c.getCache() {
		rc.cache[keyCache] = cache
	}

	rc.eventsContainer.merge(c.getEventsContainer())
}

//getConstructors exposes constructors for merge
func (rc *RuntimeContainer) getConstructors() map[string]Constructor {
	return rc.constructors
}

//getCache exposes cache for merge
func (rc *RuntimeContainer) getCache() dependencyCache {
	return rc.cache
}

//getEventsContainer exposes events for merge
func (rc *RuntimeContainer) getEventsContainer() EventsContainer {
	return *rc.eventsContainer
}
