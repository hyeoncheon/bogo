package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"prober/checks"
	"prober/exporters"

	getopt "github.com/pborman/getopt/v2"
)

func main() {
	getopt.Parse()
	targets := getopt.Args()
	run(targets)
}

func run(targets []string) {
	out := make(chan checks.PingMessage)
	exporterLock := make(chan int)

	exporter := &exporters.StdoutExporter{}
	exporter.Initialize(out, exporterLock)

	for _, t := range targets {
		go func(t string) {
			defer func(t string) {
				v := recover()
				fmt.Printf("panic on workder for %v! interruptted? %v\n", t, v)
			}(t)

			for {
				checks.Ping(t, out)
				time.Sleep(5 * time.Second)
			}
		}(t)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	for {
		s := <-sig
		fmt.Println("signal caught:", s)
		switch s {
		case syscall.SIGINT:
			fmt.Println("interrupted!")
			close(out)
		}
		break
	}

	// wait until exporter exit
	<-exporterLock
}
