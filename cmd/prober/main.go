package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"prober"
	"prober/checks"
	"prober/exporters"

	getopt "github.com/pborman/getopt/v2"
)

func main() {
	getopt.Parse()
	targets := getopt.Args()

	if len(targets) == 0 {
		prober.Err("no 'targets' from command line. checking metadata...")
		c := prober.NewMetadataClient()
		targetList, err := c.InstanceAttributeValue("targets")
		if err != nil || len(targetList) == 0 {
			prober.Err("no 'targets' in instance attributes. checking project level...")
			targetList, err = c.ProjectAttributeValue("targets")
			if err != nil || len(targetList) == 0 {
				prober.Err("no 'targets' in project attributes. how can we do...")

				prober.Err("no targets specified. abort!")
				os.Exit(0)
			}
		}

		for _, t := range strings.Split(targetList, ",") {
			targets = append(targets, strings.TrimSpace(t))
		}
	}

	prober.Info("starting prober for %v", targets)
	run(targets)
}

func run(targets []string) {
	out := make(chan prober.PingMessage)
	exporterLock := make(chan int)

	// exporter := &exporters.StdoutExporter{}
	exporter := &exporters.StackdriverExporter{}
	exporter.Initialize(out, exporterLock)

	for _, t := range targets {
		go func(t string) {
			defer func(t string) {
				v := recover()
				prober.Err("panic on workder for %v! interruptted? %v", t, v)
			}(t)

			for {
				checks.Ping(t, out) // it takes 5 to 10 secs
			}
		}(t)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	for {
		s := <-sig
		prober.Err("signal caught: %v", s)
		switch s {
		case syscall.SIGINT:
			prober.Err("interrupted!")
			close(out)
		}
		break
	}

	// wait until exporter exit
	<-exporterLock
}
