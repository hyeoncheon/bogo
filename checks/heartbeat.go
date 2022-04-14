package checks

import (
	"time"

	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	heartbeatChecker         = "heartbeat"
	heartbeatCheckerInterval = 1 * time.Minute
)

func (x *Checker) Heartbeat() error {
	x.Name = heartbeatChecker
	x.Run = heartbeatRunner
	return nil
}

func heartbeatRunner(c common.Context, opts common.PluginOptions, out chan interface{}) error {
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
		logger.Info(heartbeatChecker, " done.")
	}()
	return nil
}
