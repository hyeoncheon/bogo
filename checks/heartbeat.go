package checks

import (
	"fmt"
	"time"

	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	heartbeatChecker            = "heartbeat"
	heartbeatCheckerIntervalSec = 60
)

// RegisterHeartbeat returns a new Checker instance and it is used by StartAll().
func (*Checker) RegisterHeartbeat() *Checker {
	return &Checker{
		name:    heartbeatChecker,
		runFunc: heartbeatRunner,
	}
}

// heartbeatRunner is a Runner function for the HeartbeatChecker, the sample
// checker implementation.
// It starts goroutine that report periodic heartbeat and returns the error
// status.
func heartbeatRunner(c common.Context, opts common.PluginOptions, out chan interface{}) error {
	logger := c.Logger().WithField("checker", heartbeatChecker)

	interval, err := opts.GetIntegerOr("interval", heartbeatCheckerIntervalSec)
	if err != nil {
		return fmt.Errorf("%w: interval", common.ErrInvalidOptionValue)
	}

	c.WG().Add(1)
	go func() { // nolint
		defer c.WG().Done()

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		logger.Infof("%s every %ds", heartbeatChecker, interval)
	infinite:
		for {
			select {
			case <-c.Done():
				break infinite
			case <-ticker.C:
				out <- "heartbeat"
			}
		}
		logger.Infof("%s checker exited", heartbeatChecker)
	}()

	return nil
}
