package meari

import (
	"context"
	"net/http"

	"github.com/hyeoncheon/bogo/internal/defaults"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// defaultServer is a "echo" based simple webserver that implements the
// Server interface.
type defaultServer struct {
	*echo.Echo
	address string
}

var _ Server = &defaultServer{}

// NewDefaultServer initializes an instance of defaultServer, registers
// the "root" handler, and returns it as a Server. Currently, Server works as
// a singleton.
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

// Address implements the Server interface.
func (s *defaultServer) Address() string {
	return s.address
}

// Serve implements the Server interface.
func (s *defaultServer) Serve() error {
	return s.Echo.Start(s.address)
}

// Shutdown implements the Server interface.
func (s *defaultServer) Shutdown(c context.Context) error {
	return s.Echo.Shutdown(c)
}

// GET implements the Server interface.
func (s *defaultServer) GET(path string, handler echo.HandlerFunc) {
	s.Echo.GET(path, handler)
}
