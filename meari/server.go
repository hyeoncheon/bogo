// Package meari provides a simple webserver feature to enable interactive
// actions such as MTR looking glass.
package meari

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/hyeoncheon/bogo/handlers"
	"github.com/hyeoncheon/bogo/internal/common"
)

var errUnsupportedMethod = errors.New("unsupported method")

// for the singleton server.
var (
	serverOnce sync.Once
	server     Server
	serverErr  error // nolint
)

// Options is a basic structure that contains all options for the Server
// implementation.
type Options struct {
	Logger  common.Logger
	Address string
}

// Server is an interface for the built-in webserver.
type Server interface {
	// Address returns the address on which the server is listening.
	Address() string
	// Serve starts the webserver and wait until the server stops.
	// It has the same behavior of `(http.Server).Serve()` and always returns
	// non-nil error. Expected error is `http.ErrServerClosed`.
	Serve() error
	// Shutdown stops the webserver gracefully.
	// It is equivalent to `(http.Server).Shutdown()`
	Shutdown(context context.Context) error
	// GET registers a new GET route for the given path with the given request handler.
	GET(path string, handler echo.HandlerFunc)
}

// NewServer initializes a new `defaultServer`, registers all available
// request handlers, then returns the Server instance.
func NewServer(c common.Context, opts *common.Options) (Server, error) {
	logger := c.Logger().WithField("module", "web")

	serverOpts := &Options{
		Logger:  logger,
		Address: opts.Address,
	}

	serverOnce.Do(func() {
		server = NewDefaultServer(serverOpts)
		serverErr = nil

		for p, handler := range handlers.AllHandlers() {
			switch handler.Method {
			case http.MethodGet:
				logger.Debugf("register handler for 'GET %v'...", p)
				server.GET(p, handler.Handler)
			default:
				serverErr = fmt.Errorf("%w: %v", errUnsupportedMethod, handler)
			}
		}
	})

	return server, serverErr
}
