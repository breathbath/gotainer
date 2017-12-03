package mocks

type ConfigProvider struct{}

func (this ConfigProvider) GetItems() map[string]interface{} {
	return map[string]interface{}{
		"EnableLogging": true,
	}
}
