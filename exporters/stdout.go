package exporters

import (
	"prober"
)

type StdoutExporter struct {
}

func (e *StdoutExporter) Initialize(in chan prober.PingMessage, wait chan int) {
	prober.Info("stdout exporter: initialize exporter...")
	go e.run(in, wait)
}

func (e *StdoutExporter) run(in chan prober.PingMessage, wait chan int) {
	defer prober.Info("stdout: bye")

	for {
		s, ok := <-in
		prober.Info("stdout: got a input %v, %v", s, ok)
		if !ok {
			wait <- 1
			return
		}
	}
}
