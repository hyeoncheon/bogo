package checks

import (
	"time"

	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	heartbeatChecker         = "heartbeat"
	heartbeatCheckerInterval = 1 * time.Minute
)

func (*Checker) RegisterHeartbeat() *Checker {
	return &Checker{
		name:    heartbeatChecker,
		runFunc: heartbeatRunner,
	}
}

func heartbeatRunner(c common.Context, _ common.PluginOptions, out chan interface{}) error {
	logger := c.Logger().WithField("checker", heartbeatChecker)
	c.WG().Add(1)
	go func() {
		defer c.WG().Done()

		ticker := time.NewTicker(heartbeatCheckerInterval)
		defer ticker.Stop()

		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("panic: %v", r)
			}
		}()

	infinite:
		for {
			select {
			case <-c.Done():
				break infinite
			case <-ticker.C:
				out <- "heartbeat"
			case <-time.After(checkSleep):
			}
		}
		logger.Infof("%s checker exited", heartbeatChecker)
	}()
	return nil
}
