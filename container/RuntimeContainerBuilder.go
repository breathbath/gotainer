package container

type RuntimeContainerBuilder struct{}

func (rc RuntimeContainerBuilder) BuildContainer(trees ...Tree) Container {
	runtimeContainer := NewRuntimeContainer()

	for _, tree := range trees {
	}
}
