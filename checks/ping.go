package checks

import (
	"fmt"
	"time"

	"github.com/hyeoncheon/bogo"

	"github.com/go-ping/ping"
)

const (
	pingChecker         = "ping"
	pingCheckerInterval = 30 * time.Second
	pingInterval        = 1 * time.Second
	pingCount           = 10
)

func (x *Checker) Ping() error {
	x.Name = pingChecker
	x.Run = pingRunner
	return nil
}

func pingRunner(c bogo.Context, out chan interface{}) error {
	logger := c.Logger().WithField("checker", pingChecker)

	targets := c.GetCloudMeta("targets")

	for _, h := range targets {
		c.WG().Add(1)
		go func(host string) {
			defer c.WG().Done()
			ticker := time.NewTicker(pingCheckerInterval)
			defer ticker.Stop()

			logger.Infof("%v checker for %v started.", pingChecker, host)
		infinit:
			for {
				select {
				case <-c.Done():
					break infinit
				case <-ticker.C:
					if err := doPing(host, out); err != nil {
						if err.Error() == "panic send on closed channel" {
							logger.Warn(err)
						} else {
							logger.Error(err)
						}
					}
				case <-time.After(checkSleep):
				}
			}
			logger.Infof("%v checker for %v finished.", pingChecker, host)
		}(h)
	}
	return nil
}

func doPing(target string, out chan interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	pinger, err := ping.NewPinger(target)
	if err != nil {
		return err
	}

	pinger.Count = pingCount
	pinger.Interval = pingInterval
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
