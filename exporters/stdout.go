package exporters

import (
	"github.com/hyeoncheon/bogo"
)

type StdoutExporter struct {
}

func (e *StdoutExporter) Initialize(in chan bogo.PingMessage, wait chan int) {
	bogo.Info("stdout exporter: initialize exporter...")
	go e.run(in, wait)
}

func (e *StdoutExporter) run(in chan bogo.PingMessage, wait chan int) {
	defer bogo.Info("stdout: bye")

	for {
		s, ok := <-in
		bogo.Info("stdout: got a input %v, %v", s, ok)
		if !ok {
			wait <- 1
			return
		}
	}
}
