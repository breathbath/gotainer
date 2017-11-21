package container

type RuntimeContainerBuilder struct{}

func (rc RuntimeContainerBuilder) BuildContainerFromConfig(trees ...Tree) Container {
	runtimeContainer := NewRuntimeContainer()

	for _, tree := range trees {
		runtimeContainer = rc.buildContainer(tree)
	}

	return runtimeContainer
}

func (rc RuntimeContainerBuilder) buildContainer(tree Tree) *RuntimeContainer {
	runtimeContainer := NewRuntimeContainer()
	for _, node := range tree {
		rc.addNode(node, runtimeContainer)
	}

	return runtimeContainer
}

func (rc RuntimeContainerBuilder) addNode(node Node, container *RuntimeContainer) {
	if node.NewFunc != nil {
		rc.addNewFunc(node.Id, node.NewFunc, node.ServiceNames, container)
	} else if node.Constr != nil {
		rc.addConstr(node.Id, node.Constr, container)
	} else if node.Ev.Service != "" {
		rc.addEvent(node.Ev.Name, node.Ev.Service, container)
	} else if node.Ob.Name != "" {
		rc.addObserver(node.Ob.Event, node.Ob.Name, node.Ob.Callback, container)
	}
}

func (rc RuntimeContainerBuilder) addNewFunc(serviceId string, newFunc interface{}, serviceNames []string, container *RuntimeContainer) {
	container.AddNewMethod(serviceId, newFunc, serviceNames...)
}

func (rc RuntimeContainerBuilder) addConstr(serviceId string, constr Constructor, container *RuntimeContainer) {
	container.AddConstructor(serviceId, constr)
}

func (rc RuntimeContainerBuilder) addEvent(eventName, dependencyName string, container *RuntimeContainer) {
	container.RegisterDependencyEvent(eventName, dependencyName)
}

func (rc RuntimeContainerBuilder) addObserver(eventName, observerId string, callback interface{}, container *RuntimeContainer) {
	container.AddDependencyObserver(eventName, observerId, callback)
}
