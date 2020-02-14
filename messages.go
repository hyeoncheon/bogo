package prober

import (
	"net"
	"time"
)

type PingMessage struct {
	Addr   string
	IPAddr *net.IPAddr
	Count  int
	Loss   float64
	MinRtt time.Duration
	MaxRtt time.Duration
	AvgRtt time.Duration
	StdDev time.Duration
}
