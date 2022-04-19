package meari

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hyeoncheon/bogo/handlers"
	"github.com/hyeoncheon/bogo/internal/common"
)

type Options struct {
	Logger  common.Logger
	Address string
}

type Server interface {
	Address() string
	Start() error
	Shutdown(c context.Context) error
	GET(string, echo.HandlerFunc)
}

func NewServer(c common.Context, opts *common.Options) (Server, error) {
	logger := c.Logger().WithField("module", "web")

	serverOpts := &Options{
		Logger:  logger,
		Address: opts.Address,
	}
	server := NewDefaultServer(serverOpts)
	if server == nil {
		return nil, fmt.Errorf("could not initiate the web server: %v", serverOpts)
	}

	for p, handler := range handlers.AllHandlers() {
		switch handler.Method {
		case http.MethodGet:
			logger.Debugf("register handler for 'GET %v'...", p)
			server.GET(p, handler.Handler)
		default:
			return nil, fmt.Errorf("unsupported method for %v", handler)
		}
	}

	return server, nil
}
