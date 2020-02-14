package prober

type Exproter interface {
	Initialize(in chan PingMessage, wait chan int)
	Write(target string, val int)
}
