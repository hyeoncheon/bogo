package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hyeoncheon/bogo/meari"
)

func (h *Handler) Echo() {
	h.Path = "/echo"
	h.Method = http.MethodGet
	h.Handler = EchoHandler
}

var _ echo.HandlerFunc = EchoHandler

func EchoHandler(c echo.Context) error {
	ret := meari.RequestAll(c)
	return c.String(http.StatusOK, ret)
}
