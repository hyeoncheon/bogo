package handlers

import (
	"net/http"
	"net/http/httptest"
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

	err := MtrHandler(c)
	r.NoError(err)
	r.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	r.Contains(body, "$ mtr -bznrTP53 -c5 127.0.0.1\n")
}

func TestMtrHandler_Error(t *testing.T) {
	r := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/mtr?target=", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := MtrHandler(c)
	r.NoError(err)
	r.Equal(http.StatusOK, rec.Code)
	body := rec.Body.String()
	r.Contains(body, "$ mtr -bznrTP53 -c5 \n")
	r.Contains(body, "mtr: Failed to resolve host")
}
