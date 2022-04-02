package common

import "github.com/sirupsen/logrus"

// Logger is an interface which supports logging with fields
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

// asset DefaultLogger as Logger
var _ Logger = DefaultLogger{}

// DefaultLogger based on logrus.FieldLogger
type DefaultLogger struct {
	logrus.FieldLogger
}

func (l DefaultLogger) WithField(s string, i interface{}) Logger {
	return DefaultLogger{l.FieldLogger.WithField(s, i)}
}

func (l DefaultLogger) WithFields(m map[string]interface{}) Logger {
	return DefaultLogger{l.FieldLogger.WithFields(m)}
}

func NewDefaultLogger(level string) Logger {
	l := logrus.New()
	if lvl, err := logrus.ParseLevel(level); err != nil {
		l.Warnf("unsupported log level %v. fallback to info.", level)
		l.Level = logrus.InfoLevel
	} else {
		l.Level = lvl
	}
	return DefaultLogger{l}
}
