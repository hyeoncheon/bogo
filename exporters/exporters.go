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
