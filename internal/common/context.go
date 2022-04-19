package common

import (
	"context"
	"sync"
)

// Context is the main application context for whole components
// This context will be used to control go routines.
type Context interface {
	context.Context
	Cancel()
	Channel() chan interface{}
	WG() *sync.WaitGroup
	Logger() Logger
	Meta() MetaClient
}

// asset DefautContext for Context iplemetations
var _ Context = &defaultContext{}
var _ context.Context = &defaultContext{}

// defaultContext is the default context for bogo app.
type defaultContext struct {
	context.Context
	Options
	cancel context.CancelFunc
	ch     chan interface{}
	wg     *sync.WaitGroup
	logger Logger
	meta   MetaClient
}

// NewDefaultContext creates a new DefaultContext with cancel function
// then returns it as Context.
func NewDefaultContext(opts *Options) (Context, context.CancelFunc) {
	c, cancel := context.WithCancel(context.Background())
	return &defaultContext{
		Context: c,
		Options: *opts,
		cancel:  cancel,
		ch:      make(chan interface{}, 10),
		wg:      &sync.WaitGroup{},
		logger:  NewDefaultLogger(opts.LogLevel),
		meta:    nil,
	}, cancel
}

func (c *defaultContext) Channel() chan interface{} {
	return c.ch
}

func (c *defaultContext) Cancel() {
	c.Logger().Debug("cancelling the main context...")
	c.cancel()

	c.Logger().Debug("waiting for routines: ", c.wg, "...")
	c.wg.Wait()

	c.Logger().Debug("closing communication channel...")
	close(c.ch)
}

func (c *defaultContext) WG() *sync.WaitGroup {
	return c.wg
}

func (c *defaultContext) Logger() Logger {
	return c.logger
}

// Meta returns the context's metadata client. If the meta is nil, it will
// try to create a new one. Currently only GCE is supported.
func (c *defaultContext) Meta() MetaClient {
	if c.meta == nil {
		c.meta = NewGCEMetaClient(c)
	}
	return c.meta
}
