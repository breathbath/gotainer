package container

import (
	"errors"
	"fmt"
)

type RuntimeContainer struct {
	constructors map[string]Constructor
	cache        dependencyCache
	eventsContainer *EventsContainer
}

func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{constructors: make(map[string]Constructor), cache: newDependencyCache(), eventsContainer: NewEventsContainer()}
}

func (this *RuntimeContainer) AddConstructor(id string, constructor Constructor) {
	this.constructors[id] = constructor
}

func (this *RuntimeContainer) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	this.constructors[id] = convertNewMethodToConstructor(this, typedConstructor, constructorArgumentNames)
}

func (this *RuntimeContainer) AddDependencyObserver(eventName, observerId string, observer interface{}) {
	this.eventsContainer.addDependencyObserver(eventName, observerId, observer)
}

func (this *RuntimeContainer) RegisterDependencyEvent(eventName, dependencyName string) {
	this.eventsContainer.registerDependencyEvent(eventName, dependencyName)
}

func (this *RuntimeContainer) Scan(id string, dest interface{}) {
	baseValue := this.Get(id, true)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

func (this *RuntimeContainer) ScanNonCached(id string, dest interface{}) {
	baseValue := this.Get(id, false)
	err := copySourceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

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

func (this *RuntimeContainer) Check() {
	for dependencyName, _ := range this.constructors {
		this.Get(dependencyName, false)
	}
}
