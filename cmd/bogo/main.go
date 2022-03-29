package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/checks"
	"github.com/hyeoncheon/bogo/exporters"

	getopt "github.com/pborman/getopt/v2"
)

type Options struct {
	isDebug  bool
	exporter string
}

func main() {
	showVersion := false
	showHelp := false

	opts := &Options{
		isDebug:  false,
		exporter: "stackdriver",
	}
	getopt.SetParameters("targets...")
	getopt.FlagLong(&opts.isDebug, "debug", 'd', "debugging mode")
	getopt.FlagLong(&opts.exporter, "exporter", 'e', "set exporter")
	getopt.FlagLong(&showVersion, "version", 'v', "show version")
	getopt.FlagLong(&showHelp, "help", 'h', "show help message")

	getopt.Parse()
	targets := getopt.Args()

	if showVersion {
		bogo.Info("bogo " + bogo.Version)
		return
	}
	if showHelp {
		bogo.Info("bogo " + bogo.Version)
		getopt.Usage()
		return
	}

	if len(targets) == 0 {
		bogo.Err("no 'targets' from command line. checking metadata...")
		c := bogo.NewMetadataClient()
		targetList, err := c.InstanceAttributeValue("targets")
		if err != nil || len(targetList) == 0 {
			bogo.Err("no 'targets' in instance attributes. checking project level...")
			targetList, err = c.ProjectAttributeValue("targets")
			if err != nil || len(targetList) == 0 {
				bogo.Err("no 'targets' in project attributes. how can we do...")

				bogo.Err("no targets specified. abort!")
				os.Exit(0)
			}
		}

		for _, t := range strings.Split(targetList, ",") {
			targets = append(targets, strings.TrimSpace(t))
		}
	}

	bogo.Info("starting bogo for %v", targets)
	run(opts, targets)
}

func run(opts *Options, targets []string) {
	out := make(chan bogo.PingMessage)
	exporterLock := make(chan int)

	var exporter bogo.PingExproter
	switch opts.exporter {
	case "stackdriver":
		exporter = &exporters.StackdriverExporter{}
	default:
		exporter = &exporters.StdoutExporter{}
	}
	exporter.Initialize(out, exporterLock)

	for _, t := range targets {
		go func(t string) {
			defer func(t string) {
				v := recover()
				bogo.Err("panic on workder for %v! interruptted? %v", t, v)
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
		bogo.Err("signal caught: %v", s)
		switch s {
		case syscall.SIGINT:
			bogo.Err("interrupted!")
			close(out)
		}
		break
	}

	// wait until exporter exit
	<-exporterLock
}
