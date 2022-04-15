package handlers

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

var mtrArgs = []string{"-bznrTP53", "-c5"}

func (h *Handler) Mtr() {
	h.Path = "/mtr"
	h.Method = http.MethodGet
	h.Handler = MtrHandler
}

var _ echo.HandlerFunc = MtrHandler

func MtrHandler(c echo.Context) error {
	target := c.FormValue("target")

	ret := fmt.Sprintf("$ mtr %v %v\n", strings.Join(mtrArgs, " "), target)

	cmd := exec.Command("mtr", append(mtrArgs, target)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		ret += err.Error()
	} else {
		ret += string(out)
	}

	return c.String(http.StatusOK, ret)
}
