package checks

import (
	"fmt"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"

	"github.com/go-ping/ping"
)

const (
	pingChecker            = "ping"
	pingCheckerIntervalSec = 30
	pingIntervalSec        = 1
	pingCount              = 10
)

func (*Checker) RegisterPing() *Checker {
	return &Checker{
		name:    pingChecker,
		runFunc: pingRunner,
	}
}

func pingRunner(c common.Context, opts common.PluginOptions, out chan interface{}) error {
	logger := c.Logger().WithField("checker", pingChecker)
	logger.Debug("ping opts: ", opts)

	// setup ping parameters
	targets := opts.GetValuesOr("targets", []string{})
	if len(targets) < 1 && c.Meta() != nil {
		targets = c.Meta().AttributeValues("targets")
	}
	if len(targets) < 1 {
		return fmt.Errorf("no targets specified")
	}

	for _, t := range targets {
		if len(t) < 1 {
			return fmt.Errorf("target string should not be empty: %v", targets)
		}
	}

	checkInterval, err := opts.GetIntegerOr("check_interval", pingCheckerIntervalSec)
	if err != nil {
		return fmt.Errorf("invalid check_interval value")
	}

	// spawn ping workers for each target
	for _, h := range targets {
		c.WG().Add(1)
		go func(host string) {
			defer c.WG().Done()

			ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
			defer ticker.Stop()

			logger.Infof("%s %s every %ds", pingChecker, host, checkInterval)
		infinite:
			for {
				select {
				case <-c.Done():
					break infinite
				case <-ticker.C:
					if err := doPing(host, out); err != nil {
						if err.Error() == "panic: send on closed channel" {
							logger.Warn(err)
						} else {
							logger.Error(err)
						}
					}
				case <-time.After(checkSleep):
				}
			}
			logger.Infof("%s checker for %s exited", pingChecker, host)
		}(h)
	}
	return nil
}

// doPing runs a single turn of ping test for the given target with fixed
// configuration, then send the result to the given channel.
func doPing(target string, out chan interface{}) (err error) {
	defer func() {
		// NOTE: mainly for "send on closed channel"
		// could it be prevented by closing the channel after all checkers?
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	pinger, err := ping.NewPinger(target)
	if err != nil {
		return err
	}

	pinger.Count = pingCount
	pinger.Interval = pingIntervalSec * time.Second
	pinger.Timeout = pinger.Interval*time.Duration(pingCount) + time.Second

	if err := pinger.Run(); err != nil {
		return err
	}
	stats := pinger.Statistics()

	out <- bogo.PingMessage{
		Addr:   stats.Addr,
		IPAddr: stats.IPAddr,
		Count:  stats.PacketsSent,
		Loss:   stats.PacketLoss,
		MinRtt: stats.MinRtt,
		MaxRtt: stats.MaxRtt,
		AvgRtt: stats.AvgRtt,
		StdDev: stats.StdDevRtt,
	}
	return nil
}
