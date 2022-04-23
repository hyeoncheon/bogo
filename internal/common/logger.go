package common

import "github.com/sirupsen/logrus"

// Logger is an interface which supports logging with fields.
type Logger interface {
	// basic logging functions
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
	// extended functions
	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
}

var _ Logger = defaultLogger{}

// defaultLogger based on logrus.FieldLogger.
type defaultLogger struct {
	logrus.FieldLogger
}

// WithField returns a new Logger that has the given field, derived from the
// parent logger.
func (l defaultLogger) WithField(key string, value interface{}) Logger {
	return defaultLogger{l.FieldLogger.WithField(key, value)}
}

// WithFields returns a new Logger that has the given fields, derived from the
// parent logger.
func (l defaultLogger) WithFields(m map[string]interface{}) Logger {
	return defaultLogger{l.FieldLogger.WithFields(m)}
}

// NewDefaultLogger returns a new logrus based default logger.
func NewDefaultLogger(level string) Logger {
	l := logrus.New()
	if lvl, err := logrus.ParseLevel(level); err != nil {
		l.Warnf("unsupported log level %v. fallback to info.", level)
		l.Level = logrus.InfoLevel
	} else {
		l.Level = lvl
	}

	return defaultLogger{l}
}
