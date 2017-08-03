package mux

// logger is an interface that should be capable of satisfying
// most common log interfaces? This lets the mux log what happens
// if a logger is provided by the consumer.
type logger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Debug(string, ...interface{})
	Error(string, ...interface{})
}

const (
	infoLevel = iota
	warnLevel
	debugLevel
	errorLevel
)

// log is a wrapper around the logging functionality, this provides a common
// place for mux to attempt logging and return if the consumer has not defined
// a logger or log if a logger is provided
func (m *Mux) log(level int, format string, data ...interface{}) {
	if m.logger == nil {
		return
	}

	switch level {
	case infoLevel:
		m.logger.Info(format, data...)
		break
	case warnLevel:
		m.logger.Warn(format, data...)
		break
	case debugLevel:
		m.logger.Debug(format, data...)
		break
	case errorLevel:
		m.logger.Error(format, data...)
		break
	}
}
