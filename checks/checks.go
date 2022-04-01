package checks

import (
	"reflect"
	"time"

	"github.com/hyeoncheon/bogo"
)

const (
	checkSleep = 100 * time.Millisecond
)

type CheckRunner func(bogo.Context, chan interface{}) error

type Checker struct {
	Name string
	Run  CheckRunner
}

type checkers map[string]*Checker

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
