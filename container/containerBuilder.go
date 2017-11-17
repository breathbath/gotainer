package container

type ContainerBuilder interface {
	BuildContainer(trees ...Tree) Container
}
