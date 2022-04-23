package handlers

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

// TODO: using TCP/53 by default is not a good choice. need to be fixed.
var mtrArgs = []string{"-bznrTP53", "-c5"}

// Mtr configures the handler information for "/mtr" and is used by AllHandlers.
func (h *Handler) Mtr() {
	h.Path = "/mtr"
	h.Method = http.MethodGet
	h.Handler = MtrHandler
}

var _ echo.HandlerFunc = MtrHandler

// cmdCombinedOutput ((*exec.Cmd).CombinedOutput) runs the command and returns
// its combined standard output and standard error.
var cmdCombinedOutput = (*exec.Cmd).CombinedOutput

// MtrHandler is an echo request handler to serve MTR looking glass.
func MtrHandler(c echo.Context) error {
	target := c.FormValue("target")

	ret := fmt.Sprintf("$ mtr %v %v\n", strings.Join(mtrArgs, " "), target)

	/* #nosec G204 */
	cmd := exec.Command("mtr", append(mtrArgs, target)...)
	if out, err := cmdCombinedOutput(cmd); err != nil {
		ret += "\nexecution error: " + err.Error()
	} else {
		ret += string(out)
	}

	return c.String(http.StatusOK, ret)
}
