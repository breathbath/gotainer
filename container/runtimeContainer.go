package container

import (
	"fmt"
	"strings"
)

//RuntimeContainer creates Services at runtime with registered callbacks
type RuntimeContainer struct {
	constructors      map[string]Constructor
	cache             dependencyCache
	eventsContainer   *EventsContainer
	garbageCollectors map[string]GarbageCollectorFunc
}

//NewRuntimeContainer creates container
func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{
		constructors:      make(map[string]Constructor),
		cache:             newDependencyCache(),
		eventsContainer:   NewEventsContainer(),
		garbageCollectors: make(map[string]GarbageCollectorFunc),
	}
}

//AddConstructor registers a Callback to create a Service identified by id
func (rc *RuntimeContainer) AddConstructor(id string, constructor Constructor) {
	rc.constructors[id] = constructor
}

//AddNewMethod converts a New Service method to a valid Callback Constr
func (rc *RuntimeContainer) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	rc.constructors[id] = convertNewMethodToConstructor(rc, typedConstructor, constructorArgumentNames, id)
}

//AddDependencyObserver registers Service that will receive Config it is interested in
func (rc *RuntimeContainer) AddDependencyObserver(eventName, observerId string, observer interface{}) {
	rc.eventsContainer.addDependencyObserver(eventName, observerId, observer)
}

//RegisterDependencyEvent notifies observers about added Config
func (rc *RuntimeContainer) RegisterDependencyEvent(eventName, dependencyName string) {
	rc.eventsContainer.registerDependencyEvent(eventName, dependencyName)
}

//Scan copies a Service identified by id into a typed destination (its a pointer reference)
func (rc *RuntimeContainer) Scan(id string, dest interface{}) {
	baseValue := rc.Get(id, true)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

//ScanNonCached creates a Service every time rc method is called
func (rc *RuntimeContainer) ScanNonCached(id string, dest interface{}) {
	baseValue := rc.Get(id, false)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

//Get fetches a Service in a return argument
func (rc *RuntimeContainer) Get(id string, isCached bool) interface{} {
	dependency, ok := rc.cache.Get(id)
	if ok && isCached {
		return dependency
	}

	constructorFunc, ok := rc.constructors[id]
	if !ok {
		panic(fmt.Errorf("Unknown dependency '%s'", id))
	}

	service, err := constructorFunc(rc)
	if err != nil {
		panic(err)
	}

	rc.eventsContainer.collectDependencyEventsForService(rc, id, service)

	rc.cache.Set(id, service)

	return service
}

//Check ensures that all runtime Config are created correctly
func (rc *RuntimeContainer) Check() {
	for dependencyName := range rc.constructors {
		rc.Get(dependencyName, false)
	}

	rc.CollectGarbage()
}

//Check ensures that all runtime Config are created correctly
func (rc *RuntimeContainer) Exists(id string) bool {
	_, exists := rc.constructors[id]
	return exists
}

//Merge allows to merge containers
func (rc *RuntimeContainer) Merge(c MergeableContainer) {
	for keyConstructor, constr := range c.getConstructors() {
		if _, ok := rc.constructors[keyConstructor]; ok {
			conflictingErrorMessage := fmt.Sprintf(
				"Cannot merge containers because of non unique Service id '%s'",
				keyConstructor,
			)
			panic(conflictingErrorMessage)
		}
		rc.constructors[keyConstructor] = constr
	}

	for keyCache, cache := range c.getCache() {
		rc.cache[keyCache] = cache
	}

	rc.eventsContainer.merge(c.getEventsContainer())
}

//AddGarbageCollectFunc registers a garbage collection function to destroy a service resources
func (rc *RuntimeContainer) AddGarbageCollectFunc(serviceName string, gcFunc GarbageCollectorFunc) {
	rc.garbageCollectors[serviceName] = gcFunc
}

//CollectGarbage will call all registered garbage collection functions and return the aggregated error result
func (rc *RuntimeContainer) CollectGarbage() error {
	errs := []string{}
	for serviceName, gcFunc := range rc.garbageCollectors {
		service := rc.Get(serviceName, true)
		err := gcFunc(service)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Garbage collection errors: %s", strings.Join(errs, ", "))
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
