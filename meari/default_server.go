package meari

import (
	"context"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DefaultServer struct {
	*echo.Echo
	Address string
}

var _ Server = &DefaultServer{}

var (
	serverOnce sync.Once
	server     Server
)

func NewDefaultServer(opts *Options) Server {
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

		s.Echo.Use(middleware.Logger())

		s.GET("/", func(c echo.Context) error {
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
