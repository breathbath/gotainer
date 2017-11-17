//Package container provides logic for declaring services and their relations, as well as a
//centralised endpoint for fetching them from the Config container.
package container

//Constructor func to return a Service or an error
type Constructor func(c Container) (interface{}, error)
