package checks

import (
	"reflect"
	"time"

	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	checkSleep = 100 * time.Millisecond
)

// Runner is a function type for plugable checkers and exporters.
type Runner func(common.Context, common.PluginOptions, chan interface{}) error

// Checker couples the name and the runner of each checker.
type Checker struct {
	Name string
	Run  Runner
}

type checkers map[string]*Checker

// Checkers is a map of registered checkers, is built by init().
var Checkers = checkers{}

func init() {
	o := reflect.TypeOf(&Checker{})
	for i := 0; i < o.NumMethod(); i++ {
		m := o.Method(i)

		x := Checker{}
		m.Func.Call([]reflect.Value{reflect.ValueOf(&x)})

		if len(x.Name) > 0 && x.Run != nil {
			Checkers[x.Name] = &x
		}
	}
}
