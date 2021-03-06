package common

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	r := require.New(t)

	opt := Options{LogLevel: "info"}
	c, ccf := NewDefaultContext(&opt)
	r.NotNil(c)
	r.NotNil(ccf)

	r.IsType((*defaultContext)(nil), c)
	r.Implements((*Context)(nil), c)
	r.Implements((*context.Context)(nil), c)

	r.IsType((context.CancelFunc)(nil), ccf)

	r.IsType((chan interface{})(nil), c.Channel())
}

func TestContextCancel(t *testing.T) {
	r := require.New(t)

	opt := Options{LogLevel: "info"}
	c, ccf := NewDefaultContext(&opt)
	r.NotNil(c)
	r.NotNil(ccf)

	r.Nil(c.Err())
	c.Cancel()
	r.NotNil(c.Err())
}

func TestContextDoubleCancel(t *testing.T) {
	r := require.New(t)

	opt := Options{LogLevel: "info"}
	c, ccf := NewDefaultContext(&opt)
	r.NotNil(c)
	r.NotNil(ccf)

	r.Nil(c.Err())
	c.Cancel()
	r.NotNil(c.Err())
	c.Cancel()
}

func TestContextWG(t *testing.T) {
	r := require.New(t)

	opt := Options{LogLevel: "info"}
	c, ccf := NewDefaultContext(&opt)
	r.NotNil(c)
	r.NotNil(ccf)

	r.IsType((*sync.WaitGroup)(nil), c.WG())
}

func TestContextLogger(t *testing.T) {
	r := require.New(t)

	opt := Options{LogLevel: "info"}
	c, ccf := NewDefaultContext(&opt)
	r.NotNil(c)
	r.NotNil(ccf)

	r.Implements((*Logger)(nil), c.Logger())
}

func TestContextMeta(t *testing.T) {
	r := require.New(t)

	opt := Options{LogLevel: "info"}
	c, ccf := NewDefaultContext(&opt)
	r.NotNil(c)
	r.NotNil(ccf)

	r.IsType((MetaClient)(nil), c.Meta())
}
