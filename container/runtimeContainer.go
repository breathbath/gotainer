package container

import (
	"errors"
	"fmt"
)

type RuntimeContainer struct {
	constructors map[string]ArgumentsConstructor
	cache        servicesCache
}

func NewRuntimeContainer() *RuntimeContainer {
	return &RuntimeContainer{make(map[string]ArgumentsConstructor), newServicesCache()}
}

func (this *RuntimeContainer) AddConstructor(id string, constructor ArgumentsConstructor) {
	this.constructors[id] = constructor
}

func (this *RuntimeContainer) AddNoArgumentsConstructor(id string, constructor NoArgumentsConstructor) {
	this.constructors[id] = func(c Container) (interface{}, error) {
		return constructor()
	}
}

func (this *RuntimeContainer) AddTypedConstructor(id string, typedConstructor interface{}, constructorArgumentNames ...string) {
	this.constructors[id] = convertTypedToUntypedConstructor(this, typedConstructor, constructorArgumentNames)
}

func (this *RuntimeContainer) GetTypedService(id string, dest interface{}) {
	baseValue := this.GetService(id)
	err := copySurceVariableToDestinationVariable(baseValue, dest, id)
	if err != nil {
		panic(err)
	}
}

func (this *RuntimeContainer) GetService(id string) interface{} {
	service, ok := this.cache.Get(id)
	if ok {
		return service
	}

	constructorFunc, ok := this.constructors[id]
	if !ok {
		errStr := fmt.Sprintf("Unknown service '%s'", id)
		panic(errors.New(errStr))
	}

	result, err := constructorFunc(this)
	if err != nil {
		panic(err)
	}

	return result
}

func (this *RuntimeContainer) Check() {
	for dependencyName, _ := range this.constructors {
		this.GetService(dependencyName)
	}
}
