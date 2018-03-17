package container

type RuntimeContainerBuilder struct{}

func (rc RuntimeContainerBuilder) BuildContainerFromConfig(trees ...Tree) Container {
	runtimeContainer := NewRuntimeContainer()

	for _, tree := range trees {
		rc.addTreeToContainer(tree, runtimeContainer)
	}

	return runtimeContainer
}

func (rc RuntimeContainerBuilder) addTreeToContainer(tree Tree, c *RuntimeContainer) {
	for _, node := range tree {
		rc.addNode(node, c)
	}
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

	if node.Parameters != nil {
		rc.addParameters(node.Parameters, container)
	}
	if node.ParamProvider != nil {
		rc.addParametersProvider(node.ParamProvider, container)
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

func (rc RuntimeContainerBuilder) addParameters(parameters map[string]interface{}, container *RuntimeContainer) {
	RegisterParameters(container, parameters)
}

func (rc RuntimeContainerBuilder) addParametersProvider(parametersProvider ParametersProvider, container *RuntimeContainer) {
	parameters := parametersProvider.GetItems()
	rc.addParameters(parameters, container)
}
