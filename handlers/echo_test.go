package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestEchoHandler(t *testing.T) {
	r := require.New(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/echo?var=val", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := EchoHandler(c)
	r.NoError(err)
	r.Equal(http.StatusOK, rec.Code)
	r.Equal("GET /echo?var=val HTTP/1.1\nHost: example.com\nContent-Type: text/plain\n\nvar: val\n", rec.Body.String())
}
