package mocks

type Logger interface {
	Log(message string)
}

type NullLogger struct{}

func (nl NullLogger) Log(message string) {

}

type InMemoryLogger struct {
	messages []string
}

func (iml *InMemoryLogger) Log(message string) {
	iml.messages = append(iml.messages, message)
}

func (iml *InMemoryLogger) GetMessages() []string {
	return iml.messages
}

func BuildLogger(isLoggerEnabled bool) Logger {
	if isLoggerEnabled {
		return &InMemoryLogger{}
	}
	return NullLogger{}
}
