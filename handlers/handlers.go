// Package handlers contains request handlers for meari webserver to provide
// interactive services.
package handlers

import (
	"reflect"
	"sync"

	"github.com/labstack/echo/v4"
)

// Handler is a structure for the handler information.
type Handler struct {
	Path    string
	Method  string
	Handler echo.HandlerFunc
}

// since the handlers is actually not dynamic and reflect is expensive,
// using `init()` or `sync.Once` could be better.
var (
	handlersOnce sync.Once
	handlers     map[string]Handler
)

// AllHandlers returns all available handlers as a map of Handler.
func AllHandlers() map[string]Handler {
	handlersOnce.Do(func() {
		handlers = map[string]Handler{}

		o := reflect.TypeOf(&Handler{})
		for i := 0; i < o.NumMethod(); i++ {
			m := o.Method(i)
			x := Handler{}
			m.Func.Call([]reflect.Value{reflect.ValueOf(&x)})
			handlers[x.Path] = x
		}
	})

	return handlers
}
