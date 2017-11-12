//Package container provides logic for declaring services and their relations, as well as a
//centralised endpoint for fetching them from the dependencies container.
package container

//Should return a service or an error
type Constructor func(c Container) (interface{}, error)
