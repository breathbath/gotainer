package container

import (
	"fmt"
	"strings"
)

//RuntimeContainer creates Services at runtime with registered callbacks
type RuntimeContainer struct {
	constructors        map[string]Constructor
	newFuncConstructors map[string]NewFuncConstructor
	cache               dependencyCache
	eventsContainer     *EventsContainer
	garbageCollectors   *GarbageCollectorFuncs
	cycleDetector       *CycleDetector
	rootDependency      string
}

//NewRuntimeContainer creates container
func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{
		constructors:        make(map[string]Constructor),
		cache:               newDependencyCache(),
		eventsContainer:     NewEventsContainer(),
		garbageCollectors:   NewGarbageCollectorFuncs(),
		newFuncConstructors: make(map[string]NewFuncConstructor),
		cycleDetector:       NewCycleDetector(),
	}
}

//AddConstructor registers a Callback to create a Service identified by id, panics if id was already declared
func (rc *RuntimeContainer) AddConstructor(id string, constructor Constructor) error {
	err := rc.assertNoDuplicates(id)
	if err != nil {
		return err
	}
	rc.SetConstructor(id, constructor)

	return nil
}

//SetConstructor adds a new service if it's not existing or overrides an existing one
func (rc *RuntimeContainer) SetConstructor(id string, constructor Constructor) {
	rc.constructors[id] = constructor
}

//AddNewMethod converts a New Service method to a valid Callback Constr, panics if id already exists
func (rc *RuntimeContainer) AddNewMethod(
	id string,
	typedConstructor interface{},
	constructorArgumentNames ...string,
) error {
	err := rc.assertNoDuplicates(id)
	if err != nil {
		return err
	}

	return rc.SetNewMethod(id, typedConstructor, constructorArgumentNames...)
}

//SetNewMethod overrides an existing service declaration or adds a new one if it doesn't exist
func (rc *RuntimeContainer) SetNewMethod(
	id string,
	typedConstructor interface{},
	constructorArgumentNames ...string,
) error {
	constrFunc, err := convertNewMethodToNewFuncConstructor(
		rc,
		typedConstructor,
		constructorArgumentNames,
		id,
	)
	if err != nil {
		return err
	}

	rc.newFuncConstructors[id] = constrFunc
	return nil
}

//AddDependencyObserver registers Service that will receive Config it is interested in
func (rc *RuntimeContainer) AddDependencyObserver(eventName, observerID string, observer interface{}) error {
	return rc.eventsContainer.addDependencyObserver(eventName, observerID, observer)
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

//ScanSecure copies a service identified by id into a typed destination (its a pointer reference) and returns error on failure
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

//GetSecure fetches a Service in a return argument and returns an error rather than panics
func (rc *RuntimeContainer) GetSecure(id string, isCached bool) (interface{}, error) {
	if rc.rootDependency == "" {
		rc.rootDependency = id
	}

	defer rc.resetCycleDetectorIfNeeded(id)

	rc.cycleDetector.VisitBeforeRecursion(id)

	if rc.cycleDetector.IsEnabled() && rc.cycleDetector.HasCycle() {
		return nil, fmt.Errorf("Detected dependencies' cycle: %s", strings.Join(rc.cycleDetector.GetCycle(), "->"))
	}

	dependency, ok := rc.cache.Get(id)
	if ok && isCached {
		rc.cycleDetector.VisitAfterRecursion(id)
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
		errorMsgSuffix := fmt.Sprintf(" [check '%s' service]", id)
		if strings.Contains(err.Error(), errorMsgSuffix) {
			errorMsgSuffix = ""
		}

		return service, fmt.Errorf("%v%s", err, errorMsgSuffix)
	}

	rc.cycleDetector.VisitAfterRecursion(id)

	err = rc.eventsContainer.collectDependencyEventsForService(rc, id, service)
	if err != nil {
		return nil, err
	}

	rc.cache.Set(id, service)

	return service, nil
}

func (rc *RuntimeContainer) resetCycleDetectorIfNeeded(curDependency string) {
	if rc.rootDependency == curDependency {
		rc.rootDependency = ""
		rc.cycleDetector.Reset()
	}
}

//Check ensures that all runtime Config are created correctly
func (rc *RuntimeContainer) Check() error {
	errs := []error{}
	var err error
	for dependencyName := range rc.constructors {
		_, err = rc.GetSecure(dependencyName, false)
		if err != nil {
			errs = append(errs, err)
		}
	}

	for dependencyName := range rc.newFuncConstructors {
		_, err = rc.GetSecure(dependencyName, false)
		if err != nil {
			errs = append(errs, err)
		}
	}

	err = rc.CollectGarbage()
	if err != nil {
		errs = append(errs, err)
	}

	return mergeErrors(errs)
}

//Exists ensures that all runtime Config are created correctly
func (rc *RuntimeContainer) Exists(id string) bool {
	_, exists := rc.constructors[id]
	return exists
}

//Merge allows to merge containers
func (rc *RuntimeContainer) Merge(c MergeableContainer) error {
	for keyConstructor, constr := range c.getConstructors() {
		if _, ok := rc.constructors[keyConstructor]; ok {
			return fmt.Errorf(
				"Cannot merge containers because of non unique Service id '%s'",
				keyConstructor,
			)
		}

		rc.constructors[keyConstructor] = constr
	}

	for keyConstructor, constr := range c.getNewFuncConstructors() {
		if _, ok := rc.newFuncConstructors[keyConstructor]; ok {
			return fmt.Errorf(
				"Cannot merge containers because of non unique Service id '%s'",
				keyConstructor,
			)
		}
		rc.newFuncConstructors[keyConstructor] = constr
	}

	for keyCache, cache := range c.getCache() {
		rc.cache[keyCache] = cache
	}

	return rc.eventsContainer.merge(c.getEventsContainer())
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

//assertNoDuplicates checks if current dependency was not already declared
func (rc *RuntimeContainer) assertNoDuplicates(id string) error {
	_, constructorExists := rc.constructors[id]
	_, newFuncExists := rc.newFuncConstructors[id]

	if constructorExists || newFuncExists {
		return fmt.Errorf("Detected duplicated dependency declaration '%s'", id)
	}

	return nil
}
