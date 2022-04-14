package meari

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/hyeoncheon/bogo/internal/common"
)

const DefaultAddress = "127.0.0.1:6090"

type Options struct {
	Logger  common.Logger
	Address string
}

type Server interface {
	Start() error
	Shutdown(c context.Context) error
	GET(string, echo.HandlerFunc)
}

type DefaultServer struct {
	*echo.Echo
	Address string
}

var _ Server = &DefaultServer{}

var (
	serverOnce sync.Once
	server     Server
)

func New(opts *Options) Server {
	serverOnce.Do(func() {
		s := &DefaultServer{
			Echo:    echo.New(),
			Address: DefaultAddress,
		}
		s.Echo.HideBanner = true
		s.Echo.Debug = true

		if opts.Address != "" {
			s.Address = opts.Address
		}

		s.GET("/", func(c echo.Context) error {
			logger := c.Logger()
			logger.Info(c.Path())
			logger.Info(c.Request().Header)
			return c.String(http.StatusOK, "Hey, Bulldog!")
		})
		server = s
	})
	return server
}

func (s *DefaultServer) Start() error {
	return s.Echo.Start(s.Address)
}

func (s *DefaultServer) Shutdown(c context.Context) error {
	return s.Echo.Shutdown(c)
}

func (s *DefaultServer) GET(path string, handler echo.HandlerFunc) {
	s.Echo.GET(path, handler)
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
