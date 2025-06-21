package logging

// Logger defines the logging interface used throughout the client
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// NoopLogger implements Logger but does nothing
type NoopLogger struct{}

func (NoopLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (NoopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (NoopLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (NoopLogger) Error(msg string, keysAndValues ...interface{}) {}
