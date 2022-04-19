package meari

import (
	"context"
	"net/http"
	"sync"

	"github.com/hyeoncheon/bogo/internal/defaults"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type defaultServer struct {
	*echo.Echo
	address string
}

var _ Server = &defaultServer{}

var (
	serverOnce sync.Once
	server     Server
)

func NewDefaultServer(opts *Options) Server {
	serverOnce.Do(func() {
		s := &defaultServer{
			Echo:    echo.New(),
			address: defaults.ServerAddress,
		}
		s.Echo.HideBanner = true
		s.Echo.Debug = true

		if opts.Address != "" {
			s.address = opts.Address
		}

		s.Echo.Use(middleware.Logger())

		s.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hey, Bulldog!")
		})
		server = s
	})
	return server
}

func (s *defaultServer) Address() string {
	return s.address
}

func (s *defaultServer) Start() error {
	return s.Echo.Start(s.address)
}

func (s *defaultServer) Shutdown(c context.Context) error {
	return s.Echo.Shutdown(c)
}

func (s *defaultServer) GET(path string, handler echo.HandlerFunc) {
	s.Echo.GET(path, handler)
}
