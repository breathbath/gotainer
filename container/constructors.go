package container

type Constructor func(c Container) (interface{}, error)