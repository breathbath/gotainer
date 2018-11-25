package mocks

//Logger logs messages
type Logger interface {
	Log(message string)
}

//NullLogger null object implementation
type NullLogger struct{}

//Log doens't log anything but implements the interface
func (nl NullLogger) Log(message string) {

}

//InMemoryLogger logs messages in memory
type InMemoryLogger struct {
	messages []string
}

//Log saves message in memory
func (iml *InMemoryLogger) Log(message string) {
	iml.messages = append(iml.messages, message)
}

//GetMessages gives logged messages
func (iml *InMemoryLogger) GetMessages() []string {
	return iml.messages
}

//BuildLogger logger builder
func BuildLogger(isLoggerEnabled bool) Logger {
	if isLoggerEnabled {
		return &InMemoryLogger{}
	}
	return NullLogger{}
}
