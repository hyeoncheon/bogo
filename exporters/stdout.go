package exporters

import (
	"fmt"
	"prober/checks"
)

type StdoutExporter struct {
}

func (e *StdoutExporter) Initialize(in chan checks.PingMessage, wait chan int) {
	fmt.Printf("stdout exporter: initialize exporter...\n")
	go e.run(in, wait)
}

func (e *StdoutExporter) run(in chan checks.PingMessage, wait chan int) {
	defer fmt.Printf("stdout: bye\n")

	for {
		s, ok := <-in
		fmt.Printf("stdout: got a input %v, %v\n", s, ok)
		if !ok {
			wait <- 1
			return
		}
	}
}
