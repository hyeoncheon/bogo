package checks

import (
	"time"

	"github.com/hyeoncheon/bogo"
)

const (
	heartbeatChecker = "heartbeat"
)

func (x *Checker) Heartbeat() error {
	x.Name = heartbeatChecker
	x.Run = heartbeatRunner
	return nil
}

func heartbeatRunner(c bogo.Context, out chan interface{}) error {
	logger := c.Logger().WithField("checker", heartbeatChecker)
	c.WG().Add(1)
	go func() {
		defer c.WG().Done()
	infinit:
		for {
			select {
			case <-c.Done():
				break infinit
			case <-time.After(1 * time.Second):
				out <- "heartbeat"
			}
		}
		logger.Info(heartbeatChecker, " done.")
	}()
	return nil
}
