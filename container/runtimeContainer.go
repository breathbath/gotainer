package container

import (
	"errors"
	"fmt"
)

type RuntimeContainer struct {
	constructors map[string]Constructor
	cache        dependencyCache
}

func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{make(map[string]Constructor), newDependencyCache()}
}

func (this *RuntimeContainer) AddConstructor(id string, constructor Constructor) {
	this.constructors[id] = constructor
}

func (this *RuntimeContainer) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	this.constructors[id] = convertNewMethodToConstructor(this, typedConstructor, constructorArgumentNames)
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

	this.cache.Set(id, result)

	return result
}

func (this *RuntimeContainer) Check() {
	for dependencyName, _ := range this.constructors {
		this.Get(dependencyName, false)
	}
}
