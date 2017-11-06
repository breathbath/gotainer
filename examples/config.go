package examples

type Config struct {
	fakeDbConnectionString string
	staticFilesUrl string
}

func NewConfig() Config {
	return Config{"someConnectionString", "http://static.me/"}
}

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
