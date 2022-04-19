package common

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/defaults"
	"github.com/stretchr/testify/require"
)

func TestDefaultOptions(t *testing.T) {
	r := require.New(t)
	opts := DefaultOptions()
	r.Equal(false, opts.IsDebug)
	r.Equal("info", opts.LogLevel)
	r.Equal([]string{}, opts.Checkers)
	r.Equal([]string{"stackdriver"}, opts.Exporters)
	r.Equal(defaults.ServerAddress, opts.Address)
	r.Nil(nil)
}
