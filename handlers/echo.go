package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Echo configures the handler information for "/echo" which is used by
// AllHandlers.
func (h *Handler) Echo() {
	h.Path = "/echo"
	h.Method = http.MethodGet
	h.Handler = EchoHandler
}

var _ echo.HandlerFunc = EchoHandler

// EchoHandler is an echo request handler to echo the request from clients
// back to the client as a response body.
func EchoHandler(c echo.Context) error {
	return c.String(http.StatusOK, requestInfo(c))
}

// requestHeader returns reconstructed request headers from the given context
// as a form of string.
func requestHeader(c echo.Context) string {
	ret := fmt.Sprintf("Host: %v\n", c.Request().Host)

	for h, v := range c.Request().Header {
		ret += fmt.Sprintf("%v: %v\n", h, strings.Join(v, ", "))
	}

	return ret
}

// requestInfo returns reconstructed request message from the context.
func requestInfo(c echo.Context) string {
	r := c.Request()
	ret := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)
	ret = fmt.Sprintf("%v\n%v\n", ret, requestHeader(c))

	params, err := c.FormParams()
	if err == nil {
		for p := range params {
			ret = fmt.Sprintf("%v%v: %v\n", ret, p, c.FormValue(p))
		}
	}

	return ret
}
