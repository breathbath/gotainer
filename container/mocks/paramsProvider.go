package mocks

//ConfigProvider example struct to give some config options
type ConfigProvider struct{}

//GetItems returns the map of possible parameters
func (cp ConfigProvider) GetItems() map[string]interface{} {
	return map[string]interface{}{
		"EnableLogging": true,
	}
}
