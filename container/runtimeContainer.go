package container

import (
	"fmt"
	"strings"
)

type CycleControl struct {
	recStack       map[string]bool
	recStackSorted []string
	visited        map[string]bool
	cycleDetected  bool
	cycle          [] string
}

func (cc *CycleControl) VisitBeforeRecursion(dep string) bool {
	if cc.cycleDetected {
		return true
	}

	if isInRecStack, ok := cc.recStack[dep]; ok && isInRecStack {
		cc.registerCycle(dep)
		return true
	}

	if isInVisited, ok := cc.visited[dep]; ok && isInVisited {
		return false
	}

	cc.visited[dep] = true
	cc.recStack[dep] = true
	cc.recStackSorted = append(cc.recStackSorted, dep)

	return false
}

func (cc *CycleControl) VisitAfterRecursion(dep string) {
	cc.recStack[dep] = false
}

func (cc *CycleControl) registerCycle(dep string) {
	cc.recStackSorted = append(cc.recStackSorted, dep)
	if len(cc.recStackSorted) == 1 {
		cc.recStackSorted = append(cc.recStackSorted, cc.recStackSorted[0])
	}
	for _, cyclicDep := range cc.recStackSorted {
		if isTrue, ok := cc.recStack[cyclicDep]; ok && isTrue {
			cc.cycle = append(cc.cycle, cyclicDep)
		}
	}
	cc.cycleDetected = true
}

func (cc *CycleControl) GetCycle() []string {
	return cc.cycle
}

//RuntimeContainer creates Services at runtime with registered callbacks
type RuntimeContainer struct {
	constructors        map[string]Constructor
	newFuncConstructors map[string]NewFuncConstructor
	cache               dependencyCache
	eventsContainer     *EventsContainer
	garbageCollectors   *GarbageCollectorFuncs
	visitedDependencies map[string]bool
	nestingLevel        int
	visitedPath         []string
	cycleControl        *CycleControl
}

//NewRuntimeContainer creates container
func NewRuntimeContainer() *RuntimeContainer {
	cycleControl := CycleControl{cycle: []string{}, cycleDetected: false, recStack: make(map[string]bool), visited: make(map[string]bool)}
	return &RuntimeContainer{
		constructors:        make(map[string]Constructor),
		cache:               newDependencyCache(),
		eventsContainer:     NewEventsContainer(),
		garbageCollectors:   NewGarbageCollectorFuncs(),
		newFuncConstructors: make(map[string]NewFuncConstructor),
		visitedDependencies: make(map[string]bool),
		visitedPath:         []string{},
		nestingLevel:        0,
		cycleControl:        &cycleControl,
	}
}

//AddConstructor registers a Callback to create a Service identified by id
func (rc *RuntimeContainer) AddConstructor(id string, constructor Constructor) {
	rc.constructors[id] = constructor
}

//AddNewMethod converts a New Service method to a valid Callback Constr
func (rc *RuntimeContainer) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	rc.newFuncConstructors[id] = convertNewMethodToNewFuncConstructor(rc, typedConstructor, constructorArgumentNames, id)
}

//AddDependencyObserver registers Service that will receive Config it is interested in
func (rc *RuntimeContainer) AddDependencyObserver(eventName, observerId string, observer interface{}) {
	rc.eventsContainer.addDependencyObserver(eventName, observerId, observer)
}

//RegisterDependencyEvent notifies observers about added Config
func (rc *RuntimeContainer) RegisterDependencyEvent(eventName, dependencyName string) {
	rc.eventsContainer.registerDependencyEvent(eventName, dependencyName)
}

//Scan copies a Service identified by id into a typed destination (its a pointer reference) and panics on failure
func (rc *RuntimeContainer) Scan(id string, dest interface{}) {
	err := rc.ScanSecure(id, true, dest)
	if err != nil {
		panic(err)
	}
}

//ScanNonCached creates a Service every time rc method is called and panics on failure
func (rc *RuntimeContainer) ScanNonCached(id string, dest interface{}) {
	err := rc.ScanSecure(id, false, dest)
	if err != nil {
		panic(err)
	}
}

//Scan copies a Service identified by id into a typed destination (its a pointer reference) and returns error on failure
func (rc *RuntimeContainer) ScanSecure(id string, isCached bool, dest interface{}) error {
	baseValue, err := rc.GetSecure(id, isCached)
	if err != nil {
		return err
	}

	return copySourceVariableToDestinationVariable(baseValue, dest, id)
}

//Get fetches a Service in a return argument and panics if an error happens
func (rc *RuntimeContainer) Get(id string, isCached bool) interface{} {

	dependency, err := rc.GetSecure(id, isCached)
	if err != nil {
		panic(err)
	}

	return dependency
}

//Get fetches a Service in a return argument and returns an error rather than panics
func (rc *RuntimeContainer) GetSecure(id string, isCached bool) (interface{}, error) {
	isCyclic := rc.cycleControl.VisitBeforeRecursion(id)

	if isCyclic {
		return nil, fmt.Errorf("Detected dependencies' cycle: %s", strings.Join(rc.cycleControl.GetCycle(), "->"))
	}

	dependency, ok := rc.cache.Get(id)
	if ok && isCached {
		return dependency, nil
	}

	constructorFunc, ok := rc.constructors[id]
	var service interface{}
	var err error
	if !ok {
		newFuncConstructor, ok := rc.newFuncConstructors[id]
		if !ok {
			return dependency, fmt.Errorf("Unknown dependency '%s'", id)
		}
		service, err = newFuncConstructor(rc, isCached)
	} else {
		service, err = constructorFunc(rc)
	}

	if err != nil {
		return service, fmt.Errorf("%v [check '%s' service]", err, id)
	}

	rc.cycleControl.VisitAfterRecursion(id)

	rc.eventsContainer.collectDependencyEventsForService(rc, id, service)

	rc.cache.Set(id, service)

	return service, nil
}

//Check ensures that all runtime Config are created correctly
func (rc *RuntimeContainer) Check() {
	for dependencyName := range rc.constructors {
		rc.Get(dependencyName, false)
	}

	for dependencyName := range rc.newFuncConstructors {
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

	for keyConstructor, constr := range c.getNewFuncConstructors() {
		if _, ok := rc.newFuncConstructors[keyConstructor]; ok {
			conflictingErrorMessage := fmt.Sprintf(
				"Cannot merge containers because of non unique Service id '%s'",
				keyConstructor,
			)
			panic(conflictingErrorMessage)
		}
		rc.newFuncConstructors[keyConstructor] = constr
	}

	for keyCache, cache := range c.getCache() {
		rc.cache[keyCache] = cache
	}

	rc.eventsContainer.merge(c.getEventsContainer())
}

//AddGarbageCollectFunc registers a garbage collection function to destroy a service resources
func (rc *RuntimeContainer) AddGarbageCollectFunc(serviceName string, gcFunc GarbageCollectorFunc) {
	rc.garbageCollectors.Add(serviceName, gcFunc)
}

//CollectGarbage will call all registered garbage collection functions and return the aggregated error result
func (rc *RuntimeContainer) CollectGarbage() error {
	errs := []string{}
	rc.garbageCollectors.Range(func(gcName string, gcFunc GarbageCollectorFunc) bool {
		service, err := rc.GetSecure(gcName, true)
		if err != nil {
			errs = append(errs, err.Error())
		}

		err = gcFunc(service)
		if err != nil {
			errs = append(errs, err.Error())
		}

		return true
	})

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Garbage collection errors: %s", strings.Join(errs, ", "))
}

//getConstructors exposes constructors for merge
func (rc *RuntimeContainer) getConstructors() map[string]Constructor {
	return rc.constructors
}

//getConstructors exposes new func constructors for merge
func (rc *RuntimeContainer) getNewFuncConstructors() map[string]NewFuncConstructor {
	return rc.newFuncConstructors
}

//getCache exposes cache for merge
func (rc *RuntimeContainer) getCache() dependencyCache {
	return rc.cache
}

//getEventsContainer exposes events for merge
func (rc *RuntimeContainer) getEventsContainer() EventsContainer {
	return *rc.eventsContainer
}
