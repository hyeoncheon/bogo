package checks

import (
	"reflect"
	"time"

	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	checkSleep = 100 * time.Millisecond
)

// Checker couples the name and the runner of each checker.
type Checker struct {
	name    string
	runFunc common.Runner
}

var _ common.Plugin = &Checker{}

// Name implements common.Plugin
func (p *Checker) Name() string {
	return p.name
}

// Run implements common.Plugin
func (p *Checker) Run(c common.Context, opts common.PluginOptions, ch chan interface{}) error {
	return p.runFunc(c, opts, ch)
}

func StartAll(c common.Context, opts *common.Options, ch chan interface{}) {
	logger := c.Logger().WithField("module", "checker")

	for _, x := range common.Plugins(reflect.TypeOf(&Checker{})) {
		x := x.(common.Plugin)
		if len(opts.Checkers) > 0 && !common.Contains(opts.Checkers, x.Name()) {
			logger.Debugf("%v is not on the checker list. skipping...", x.Name())
			continue
		}
		copts := opts.CheckerOptions[x.Name()]
		logger.Debugf("--- checker: %s %v with %v", x.Name(), x, copts)
		logger.Infof("starting checker %v...", x.Name())
		if err := x.Run(c, copts, ch); err != nil {
			logger.Errorf("%s checker was aborted: %v", x.Name(), err)
			// TODO: should returns error?
		}
	}
}
