package examples

//Config simulates some configuration for your app
type Config struct {
	fakeDbConnectionString string
	staticFilesUrl         string
}

//NewConfig Config Constructor
func NewConfig() Config {
	return Config{"someConnectionString", "http://static.me/"}
}

//GetValue returns individual config options
func (config Config) GetValue(key string) string {
	switch key {
	case "fakeDbConnectionString":
		return config.fakeDbConnectionString
	case "staticFilesUrl":
		return config.staticFilesUrl
	default:
		return ""
	}
}
