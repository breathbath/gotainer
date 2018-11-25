package container

//GarbageCollectorFunc typed func to recognise user defined garbage collection funcs
type GarbageCollectorFunc func(service interface{}) error
