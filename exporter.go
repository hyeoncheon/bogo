package prober

import "prober/checks"

type Exproter interface {
	Initialize(in chan checks.PingMessage, wait chan int)
	Write(target string, val int)
}
