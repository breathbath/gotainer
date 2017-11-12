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

func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{constructors: make(map[string]Constructor), cache: newDependencyCache(), eventsContainer: NewEventsContainer()}
}

//AddConstructor registers a callback to create a service identified by id
func (this *RuntimeContainer) AddConstructor(id string, constructor Constructor) {
	this.constructors[id] = constructor
}

//AddNewMethod converts a New service method to a valid callback constructor
func (this *RuntimeContainer) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	this.constructors[id] = convertNewMethodToConstructor(this, typedConstructor, constructorArgumentNames)
}

//AddDependencyObserver registers service that will receive dependencies it is interested in
func (this *RuntimeContainer) AddDependencyObserver(eventName, observerId string, observer interface{}) {
	this.eventsContainer.addDependencyObserver(eventName, observerId, observer)
}

//RegisterDependencyEvent notifies observers about added dependencies
func (this *RuntimeContainer) RegisterDependencyEvent(eventName, dependencyName string) {
	this.eventsContainer.registerDependencyEvent(eventName, dependencyName)
}

//Scan copies a service identified by id into a typed destination (its a pointer reference)
func (this *RuntimeContainer) Scan(id string, dest interface{}) {
	baseValue := this.Get(id, true)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

//ScanNonCached creates a service every time this method is called
func (this *RuntimeContainer) ScanNonCached(id string, dest interface{}) {
	baseValue := this.Get(id, false)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

//Get fetches a service in a return argument
func (this *RuntimeContainer) Get(id string, isCached bool) interface{} {
	dependency, ok := this.cache.Get(id)
	if ok && isCached {
		return dependency
	}

	constructorFunc, ok := this.constructors[id]
	if !ok {
		errStr := fmt.Sprintf("Unknown dependency '%s'", id)
		panic(errors.New(errStr))
	}

	result, err := constructorFunc(this)
	if err != nil {
		panic(err)
	}

	this.eventsContainer.notifyObserverAboutDependency(*this, id, result)

	this.cache.Set(id, result)

	return result
}

//Check ensures that all runtime dependencies are created correctly
func (this *RuntimeContainer) Check() {
	for dependencyName, _ := range this.constructors {
		this.Get(dependencyName, false)
	}
}

//Merge allows to merge containers
func (this *RuntimeContainer) Merge(c MergeableContainer) {
	for keyConstructor, constr := range c.getConstructors() {
		this.constructors[keyConstructor] = constr
	}

	for keyCache, cache := range c.getCache() {
		this.cache[keyCache] = cache
	}

	this.eventsContainer.merge(c.getEventsContainer())
}

//getConstructors exposes constructors for merge
func (this *RuntimeContainer) getConstructors() map[string]Constructor {
	return this.constructors
}

//getCache exposes cache for merge
func (this *RuntimeContainer) getCache() dependencyCache {
	return this.cache
}

//getEventsContainer exposes events for merge
func (this *RuntimeContainer) getEventsContainer() EventsContainer {
	return *this.eventsContainer
}
