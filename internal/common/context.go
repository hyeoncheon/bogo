package common

import (
	"context"
	"sync"
)

// Options is a struct to store command line options
type Options struct {
	IsDebug  bool
	LogLevel string
	// Checkers constains selected checkers. All available checkers will be
	// selected if this field is empty.
	Checkers []string
	// Exporters contains selected exporters. All available exporters will be
	// selected if this field is empty.
	Exporters []string

	CheckerOptions  map[string]PluginOptions
	ExporterOptions map[string]PluginOptions
}

// Context is the main application context for whole components
// This context will be used to control go routines.
type Context interface {
	context.Context
	Cancel()
	WG() *sync.WaitGroup
	Logger() Logger
	Meta() MetaClient
}

// asset DefautContext for Context iplemetations
var _ Context = &DefaultContext{}
var _ context.Context = &DefaultContext{}

// DefaultContext is the default context for bogo app.
type DefaultContext struct {
	context.Context
	Options
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	logger Logger
	meta   MetaClient
}

// NewDefaultContext creates a new DefaultContext with cancel function
// then returns it as Context.
func NewDefaultContext(opts Options) (Context, context.CancelFunc) {
	c, cancel := context.WithCancel(context.Background())
	return &DefaultContext{
		Context: c,
		Options: opts,
		cancel:  cancel,
		wg:      &sync.WaitGroup{},
		logger:  NewDefaultLogger(opts.LogLevel),
		meta:    nil,
	}, cancel
}

func (c *DefaultContext) Cancel() {
	c.cancel()
}

func (c *DefaultContext) WG() *sync.WaitGroup {
	return c.wg
}

func (c *DefaultContext) Logger() Logger {
	return c.logger
}

// Meta returns the context's metadata client. If the meta is nil, it will
// try to create a new one. Currently only GCE is supported.
func (c *DefaultContext) Meta() MetaClient {
	if c.meta == nil {
		c.meta = NewGCEMetaClient(c)
	}
	return c.meta
}
