package container

//RuntimeContainerBuilder builds a Runtime container
type RuntimeContainerBuilder struct{}

//BuildContainerFromConfig given a config it will build a container, panics if config is wrong
func (rc RuntimeContainerBuilder) BuildContainerFromConfig(trees ...Tree) (Container, error) {
	runtimeContainer, err := rc.BuildContainerFromConfigSecure(trees...)

	return runtimeContainer, err
}

//BuildContainerFromConfigSecure given a config it will build a container, if config is wrong an error is returned
func (rc RuntimeContainerBuilder) BuildContainerFromConfigSecure(trees ...Tree) (Container, error) {
	runtimeContainer := NewRuntimeContainer()

	mergedTree := rc.mergeTrees(trees)

	err := ValidateConfigSecure(mergedTree)
	if err != nil {
		return runtimeContainer, err
	}

	err = rc.addTreeToContainer(mergedTree, runtimeContainer)

	return runtimeContainer, err
}

func (rc RuntimeContainerBuilder) addTreeToContainer(tree Tree, c *RuntimeContainer) (err error) {
	errors := []error{}
	for _, node := range tree {
		err = rc.addNode(node, c)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return mergeErrors(errors)
}

func (rc RuntimeContainerBuilder) addNode(node Node, container *RuntimeContainer) (err error) {
	if node.NewFunc != nil {
		rc.addNewFunc(node.ID, node.NewFunc, node.ServiceNames, container)
	} else if node.Constr != nil {
		rc.addConstr(node.ID, node.Constr, container)
	} else if node.Ev.Service != "" {
		rc.addEvent(node.Ev.Name, node.Ev.Service, container)
	} else if node.Ob.Name != "" {
		rc.addObserver(node.Ob.Event, node.Ob.Name, node.Ob.Callback, container)
	}

	if node.Parameters != nil {
		err = rc.addParameters(node.Parameters, container)
		if err != nil {
			return err
		}
	}
	if node.ParamProvider != nil {
		err = rc.addParametersProvider(node.ParamProvider, container)
		if err != nil {
			return err
		}
	}

	if node.GarbageFunc != nil {
		container.AddGarbageCollectFunc(node.ID, node.GarbageFunc)
	}

	return nil
}

func (rc RuntimeContainerBuilder) addNewFunc(serviceID string, newFunc interface{}, serviceNames []string, container *RuntimeContainer) {
	container.AddNewMethod(serviceID, newFunc, serviceNames...)
}

func (rc RuntimeContainerBuilder) addConstr(serviceID string, constr Constructor, container *RuntimeContainer) {
	container.AddConstructor(serviceID, constr)
}

func (rc RuntimeContainerBuilder) addEvent(eventName, dependencyName string, container *RuntimeContainer) {
	container.RegisterDependencyEvent(eventName, dependencyName)
}

func (rc RuntimeContainerBuilder) addObserver(eventName, observerID string, callback interface{}, container *RuntimeContainer) {
	container.AddDependencyObserver(eventName, observerID, callback)
}

func (rc RuntimeContainerBuilder) addParameters(parameters map[string]interface{}, container *RuntimeContainer) error {
	return RegisterParameters(container, parameters)
}

func (rc RuntimeContainerBuilder) addParametersProvider(parametersProvider ParametersProvider, container *RuntimeContainer) error {
	parameters := parametersProvider.GetItems()
	return rc.addParameters(parameters, container)
}

func (rc RuntimeContainerBuilder) mergeTrees(trees []Tree) Tree {
	mergedTree := []Node{}
	for _, tree := range trees {
		mergedTree = append(mergedTree, tree...)
	}

	return Tree(mergedTree)
}
