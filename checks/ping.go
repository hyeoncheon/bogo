package checks

import (
	"fmt"
	"net"
	"time"

	"github.com/sparrc/go-ping"
)

func Ping(target string, out chan PingMessage) {
	fmt.Printf("ping: hello %v: %v\n", target, 9)

	pinger, err := ping.NewPinger(target)
	if err != nil {
		panic(err)
	}

	pinger.Count = 3
	pinger.Run()
	stats := pinger.Statistics()
	fmt.Println("stat:", stats.PacketLoss, stats.IPAddr, stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

	mesg := PingMessage{
		Addr:   stats.Addr,
		IPAddr: stats.IPAddr,
		Loss:   stats.PacketLoss,
		MinRtt: stats.MinRtt,
		MaxRtt: stats.MaxRtt,
		AvgRtt: stats.AvgRtt,
		StdDev: stats.StdDevRtt,
	}
	out <- mesg
}

type PingMessage struct {
	Addr   string
	IPAddr *net.IPAddr
	Loss   float64
	MinRtt time.Duration
	MaxRtt time.Duration
	AvgRtt time.Duration
	StdDev time.Duration
}
