package checks

import (
	"time"

	"github.com/hyeoncheon/bogo"

	"github.com/go-ping/ping"
)

const (
	pingChecker = "ping"
)

const checkPerMinute = 3
const count = 10
const intervalMilli = 1000
const timeoutMilli = 1000

func (x *Checker) Pinging() error {
	x.Name = pingChecker
	x.Run = pingRunner
	return nil
}

func pingRunner(c bogo.Context, out chan interface{}) error {
	logger := c.Logger().WithField("checker", pingChecker)
	c.WG().Add(1)
	go func() {
		defer c.WG().Done()
	infinit:
		for {
			select {
			case <-c.Done():
				break infinit
			case <-time.After(1 * time.Second):
				out <- "ping test"
			}
		}
		logger.Info("pingRunner done.")
	}()
	return nil
}

func Ping(target string, out chan bogo.PingMessage) {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		panic(err)
	}

	pinger.Count = count
	pinger.Interval = intervalMilli * time.Millisecond
	pinger.Timeout = time.Duration(count)*pinger.Interval + time.Second

	pinger.Run()
	stats := pinger.Statistics()
	bogo.Info("stat: %v %v %v %v %v %v %v %v",
		stats.IPAddr, stats.PacketsRecv, stats.PacketsSent, stats.PacketLoss,
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

	mesg := bogo.PingMessage{
		Addr:   stats.Addr,
		IPAddr: stats.IPAddr,
		Count:  stats.PacketsSent,
		Loss:   stats.PacketLoss,
		MinRtt: stats.MinRtt,
		MaxRtt: stats.MaxRtt,
		AvgRtt: stats.AvgRtt,
		StdDev: stats.StdDevRtt,
	}
	out <- mesg
	time.Sleep((time.Minute/checkPerMinute - time.Duration(count)*pinger.Interval))
}
