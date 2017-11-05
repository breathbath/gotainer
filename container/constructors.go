package container

type ArgumentsConstructor func(c Container) (interface{}, error)

type NoArgumentsConstructor func() (interface{}, error)