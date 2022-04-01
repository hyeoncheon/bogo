package bogo

import (
	"context"
	"sync"
)

type Options struct {
	IsDebug  bool
	LogLevel string
	Exporter string
}

type Context interface {
	context.Context
	WG() *sync.WaitGroup
	Logger() Logger
	GetCloudMeta(string) []string
}

// asset iplemetation
var _ Context = &DefaultContext{}
var _ context.Context = &DefaultContext{}

type DefaultContext struct {
	context.Context
	Options
	wg     *sync.WaitGroup
	logger Logger
}

func NewDefaultContext(opts Options) (Context, context.CancelFunc) {
	c, cancel := context.WithCancel(context.Background())
	return &DefaultContext{
		Context: c,
		Options: opts,
		wg:      &sync.WaitGroup{},
		logger:  NewDefaultLogger(opts.LogLevel),
	}, cancel
}

func (c DefaultContext) WG() *sync.WaitGroup {
	return c.wg
}

func (c *DefaultContext) Logger() Logger {
	return c.logger
}

func (c *DefaultContext) GetCloudMeta(key string) []string {
	return []string{"www.google.com", "ns.kornet.net"}
}
