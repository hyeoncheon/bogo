package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefaultLogger(t *testing.T) {
	r := require.New(t)

	l := NewDefaultLogger("debug")
	r.IsType(defaultLogger{}, l)
	r.Implements((*Logger)(nil), l)

	l = NewDefaultLogger("")
	r.IsType(defaultLogger{}, l)
	r.Implements((*Logger)(nil), l)
	l.Info("should work")
}

func TestWithField(t *testing.T) {
	r := require.New(t)

	l := NewDefaultLogger("debug")
	r.IsType(defaultLogger{}, l)
	r.Implements((*Logger)(nil), l)
	fl := l.WithField("name", "gildong")
	r.IsType(defaultLogger{}, fl)
	fl = l.WithField("age", 1)
	r.IsType(defaultLogger{}, fl)
	fl = l.WithFields(map[string]interface{}{"age": 1, "name": "gildong"})
	r.IsType(defaultLogger{}, fl)
}
