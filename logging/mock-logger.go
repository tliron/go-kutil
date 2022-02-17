package logging

var MOCK_LOGGER MockLogger

//
// MockLogger
//

type MockLogger struct{}

// Logger interface

func (self MockLogger) NewMessage(level Level, depth int) Message {
	return nil
}

func (self MockLogger) Log(level Level, depth int, message string) {
}

func (self MockLogger) Logf(level Level, depth int, format string, values ...interface{}) {
}

func (self MockLogger) Critical(message string) {
}

func (self MockLogger) Criticalf(format string, values ...interface{}) {
}

func (self MockLogger) Error(message string) {
}

func (self MockLogger) Errorf(format string, values ...interface{}) {
}

func (self MockLogger) Warning(message string) {
}

func (self MockLogger) Warningf(format string, values ...interface{}) {
}

func (self MockLogger) Notice(message string) {
}

func (self MockLogger) Noticef(format string, values ...interface{}) {
}

func (self MockLogger) Info(message string) {
}

func (self MockLogger) Infof(format string, values ...interface{}) {
}

func (self MockLogger) Debug(message string) {
}

func (self MockLogger) Debugf(format string, values ...interface{}) {
}

func (self MockLogger) AllowLevel(level Level) bool {
	return false
}

func (self MockLogger) SetMaxLevel(level Level) {
}
