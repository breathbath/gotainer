package container

//Builder interface which can be used to declare multiple different container builders
type Builder interface {
	BuildContainer(trees ...Tree) Container
}
