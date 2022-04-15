package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Echo() {
	h.Path = "/echo"
	h.Method = http.MethodGet
	h.Handler = EchoHandler
}

var _ echo.HandlerFunc = EchoHandler

func EchoHandler(c echo.Context) error {
	ret := RequestAll(c)
	return c.String(http.StatusOK, ret)
}

func RequestHeader(c echo.Context) string {
	ret := fmt.Sprintf("Host: %v\n", c.Request().Host)
	headers := c.Request().Header
	for h, v := range headers {
		ret += fmt.Sprintf("%v: %v\n", h, strings.Join(v, ", "))
	}
	return ret
}

func RequestAll(c echo.Context) string {
	r := c.Request()
	ret := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)
	ret = fmt.Sprintf("%v\n%v\n", ret, RequestHeader(c))
	params, err := c.FormParams()
	if err == nil {
		for p := range params {
			ret = fmt.Sprintf("%v%v: %v\n", ret, p, c.FormValue(p))
		}
	}
	return ret
}
