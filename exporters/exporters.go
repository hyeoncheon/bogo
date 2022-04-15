package exporters

import (
	"reflect"

	"github.com/hyeoncheon/bogo/internal/common"
)

// Exporter couples the name and the runner of each exporter.
type Exporter struct {
	Name string
	Run  common.Runner
}

type exporters map[string]*Exporter

// Exporters is a map of registered exporters, is built by init().
var Exporters = exporters{}

func init() {
	/* NOTE: DUP-2356990b70031abca66a77451c35be91 */
	o := reflect.TypeOf(&Exporter{})
	for i := 0; i < o.NumMethod(); i++ {
		m := o.Method(i)

		x := Exporter{}
		m.Func.Call([]reflect.Value{reflect.ValueOf(&x)})

		if len(x.Name) > 0 && x.Run != nil {
			Exporters[x.Name] = &x
		}
	}
}

func StartAll(c common.Context, opts *common.Options, ch chan interface{}) {
	logger := c.Logger().WithField("module", "checker")

	for k, x := range Exporters {
		if len(opts.Exporters) > 0 && !common.Contains(opts.Exporters, k) {
			logger.Debugf("%v is not on the exporter list. skipping...", k)
			continue
		}
		copts := opts.ExporterOptions[k]
		logger.Debug("--- exporter:", k, x, copts)
		logger.Info("starting exporter ", k, "...")
		x.Run(c, copts, ch)
	}
}
