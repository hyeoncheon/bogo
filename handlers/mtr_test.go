package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestMtrHandler(t *testing.T) {
	r := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/mtr?target=127.0.0.1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalFunc := cmdCombinedOutput
	defer func() {
		cmdCombinedOutput = originalFunc
	}()
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte(cmd.String()), nil
	}

	err := MtrHandler(c)
	r.NoError(err)
	r.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	r.Contains(body, "$ mtr -bznrTP53 -c5 127.0.0.1\n")
}

func TestMtrHandler_NoHost(t *testing.T) {
	r := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/mtr?target=", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalFunc := cmdCombinedOutput
	defer func() {
		cmdCombinedOutput = originalFunc
	}()
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte(cmd.String() + "\nFailed to resolve host"), nil
	}

	err := MtrHandler(c)
	r.NoError(err)
	r.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	r.Contains(body, "$ mtr -bznrTP53 -c5 \n")
	r.Contains(body, "Failed to resolve host")
}

func TestMtrHandler_Error(t *testing.T) {
	r := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/mtr?target=", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalFunc := cmdCombinedOutput
	defer func() {
		cmdCombinedOutput = originalFunc
	}()
	cmdCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		return []byte(""), errors.New("unknown")
	}

	err := MtrHandler(c)
	r.NoError(err)
	r.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	r.Contains(body, "execution error: unknown")
}
