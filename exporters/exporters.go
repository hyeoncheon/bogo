// Package exporters contains all exporters to provide data delivery to the
// external systems such as Cloud Monitoring or simply stadard out.
package exporters

import (
	"reflect"

	"github.com/hyeoncheon/bogo/internal/common"
)

// Exporter couples the name and the runner of each exporter.
type Exporter struct {
	name    string
	runFunc common.Runner
}

var _ common.Plugin = &Exporter{}

// Name implements common.Plugin.
func (p *Exporter) Name() string {
	return p.name
}

// Run implements common.Plugin.
func (p *Exporter) Run(c common.Context, opts common.PluginOptions, ch chan interface{}) error {
	return p.runFunc(c, opts, ch)
}

// StartAll executes runners for all exporters and returns number of successful
// executions.
func StartAll(c common.Context, opts *common.Options, ch chan interface{}) int {
	logger := c.Logger().WithField("module", "exporter")
	n := 0

	for _, x := range common.Plugins(reflect.TypeOf(&Exporter{})) {
		x, _ := x.(common.Plugin)
		if len(opts.Exporters) > 0 && !common.Contains(opts.Exporters, x.Name()) {
			logger.Debugf("%v is not on the exporter list. skipping...", x.Name())
			continue //nolint
		}

		eopts := opts.ExporterOptions[x.Name()]
		logger.Debugf("--- exporter: %s %v with %v", x.Name(), x, eopts)

		logger.Infof("starting exporter %v...", x.Name())

		if err := x.Run(c, eopts, ch); err != nil {
			// TODO: should returns error?
			logger.Errorf("%s exporter was aborted: %v", x.Name(), err)
		} else {
			n++
		}
	}

	return n
}
