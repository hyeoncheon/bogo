package bogo

import (
	"net"
	"time"
)

// PingMessage is a type specific message for ping statistics
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
