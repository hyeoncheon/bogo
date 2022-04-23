package checks

import (
	"errors"
	"fmt"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"

	"github.com/go-ping/ping"
)

const (
	pingChecker            = "ping"
	pingCheckerIntervalSec = 30
	pingIntervalMilliSec   = 1000
	pingCount              = 10
)

var (
	errNoTargetsSpecified           = errors.New("no targets specified")
	errTargetStringShouldNotBeEmpty = errors.New("target string should not be empty")
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

	targets, err := getTarget(c, &opts)
	if err != nil {
		return err
	}

	checkInterval, err := opts.GetIntegerOr("check_interval", pingCheckerIntervalSec)
	if err != nil {
		return fmt.Errorf("%w: check_interval", common.ErrInvalidOptionValue)
	}

	pingInterval, err := opts.GetIntegerOr("ping_interval", pingIntervalMilliSec)
	if err != nil {
		return fmt.Errorf("%w: check_interval", common.ErrInvalidOptionValue)
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
					m, err := doPing(host, pingInterval)
					if err != nil {
						logger.Error(err)
						break infinite
					}
					out <- m
				case <-time.After(checkSleep):
				}
			}
			logger.Infof("%s checker for %s exited", pingChecker, host)
		}(h)
	}
	return nil
}

func getTarget(c common.Context, opts *common.PluginOptions) ([]string, error) {
	targets := opts.GetValuesOr("targets", []string{})

	if len(targets) < 1 && c.Meta() != nil {
		targets = c.Meta().AttributeValues("targets")
	}

	if len(targets) < 1 {
		return targets, errNoTargetsSpecified
	}

	for _, t := range targets {
		if len(t) < 1 {
			return targets, fmt.Errorf("%w: %v", errTargetStringShouldNotBeEmpty, targets)
		}
	}

	return targets, nil
}

// doPing runs a single turn of ping test for the given target with fixed
// configuration, then returns the result with error status.
func doPing(target string, interval int) (bogo.PingMessage, error) {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return bogo.PingMessage{}, err
	}

	pinger.Count = pingCount
	pinger.Interval = time.Duration(interval) * time.Millisecond
	pinger.Timeout = pinger.Interval*time.Duration(pingCount) + time.Second

	if err := pinger.Run(); err != nil {
		return bogo.PingMessage{}, err
	}
	stats := pinger.Statistics()

	return bogo.PingMessage{
		Addr:   stats.Addr,
		IPAddr: stats.IPAddr,
		Count:  stats.PacketsSent,
		Loss:   stats.PacketLoss,
		MinRtt: stats.MinRtt,
		MaxRtt: stats.MaxRtt,
		AvgRtt: stats.AvgRtt,
		StdDev: stats.StdDevRtt,
	}, nil
}
