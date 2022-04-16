package handlers

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Path    string
	Method  string
	Handler echo.HandlerFunc
}

func AllHandlers() map[string]Handler {
	handlers := map[string]Handler{}

	o := reflect.TypeOf(&Handler{})
	for i := 0; i < o.NumMethod(); i++ {
		m := o.Method(i)
		x := Handler{}
		m.Func.Call([]reflect.Value{reflect.ValueOf(&x)})
		handlers[x.Path] = x
	}

	return handlers
}
