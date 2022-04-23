package common

import "errors"

// for test.
type Plugger struct {
	name    string
	runFunc Runner
}

// Name implements Plugin.
func (p *Plugger) Name() string {
	return p.name
}

// Run implements Plugin.
func (p *Plugger) Run(c Context, o PluginOptions, ch chan interface{}) error {
	return p.runFunc(c, o, ch)
}

var _ Plugin = &Plugger{}

func (*Plugger) RegisterMummy() *Plugger {
	return &Plugger{
		name: "mummy",
		runFunc: func(Context, PluginOptions, chan interface{}) error {
			return errors.New("mummy")
		},
	}
}

func (*Plugger) RegisterDummy() *Plugger {
	return &Plugger{
		name: "dummy",
		runFunc: func(Context, PluginOptions, chan interface{}) error {
			return nil
		},
	}
}

func (*Plugger) DummyFunc() *Plugger {
	return &Plugger{
		name: "dummy",
		runFunc: func(Context, PluginOptions, chan interface{}) error {
			return nil
		},
	}
}
