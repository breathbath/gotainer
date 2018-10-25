package container

type ServiceRequestMock struct {
	id       string
	isCached bool
}
type ContainerInterfaceMock struct {
	ServicesRequested []ServiceRequestMock
	service           interface{}
}

func (ci *ContainerInterfaceMock) AddConstructor(id string, constructor Constructor) {

}
func (ci *ContainerInterfaceMock) AddNewMethod(id string, typedConstructor interface{}, constructorArgumentNames ...string) {

}
func (ci *ContainerInterfaceMock) Scan(id string, dest interface{}) {

}
func (ci *ContainerInterfaceMock) ScanNonCached(id string, dest interface{}) {

}

func (ci *ContainerInterfaceMock) ScanSecure(id string, isCached bool, dest interface{}) error {
	return nil
}

func (ci *ContainerInterfaceMock) Get(id string, isCached bool) interface{} {
	ci.ServicesRequested = []ServiceRequestMock{}
	ci.ServicesRequested = append(ci.ServicesRequested, ServiceRequestMock{id: id, isCached: isCached})
	return ci.service
}

func (ci *ContainerInterfaceMock) GetSecure(id string, isCached bool) (interface{}, error) {
	return ci.Get(id, isCached), nil
}

func (ci *ContainerInterfaceMock) Check() {

}

func (ci *ContainerInterfaceMock) SetServiceToReturn(service interface{}) {
	ci.service = service
}

func (ci *ContainerInterfaceMock) Exists(id string) bool {
	return true
}

func (ci *ContainerInterfaceMock) AddGarbageCollectFunc(serviceName string, gcFunc GarbageCollectorFunc) {
}

func (ci *ContainerInterfaceMock) CollectGarbage() error {
	return nil
}
