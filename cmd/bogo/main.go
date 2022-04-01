package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/checks"

	getopt "github.com/pborman/getopt/v2"
)

func main() {
	showVersion := false
	showHelp := false

	opts := bogo.Options{
		IsDebug:  false,
		LogLevel: "info",
		Exporter: "stackdriver",
	}
	getopt.SetParameters("targets...")
	getopt.FlagLong(&opts.IsDebug, "debug", 'd', "debugging mode")
	getopt.FlagLong(&opts.Exporter, "exporter", 'e', "set exporter")
	getopt.FlagLong(&opts.LogLevel, "log", 'l', "log level. (debug, info, warn, or error)")
	getopt.FlagLong(&showVersion, "version", 'v', "show version")
	getopt.FlagLong(&showHelp, "help", 'h', "show help message")

	getopt.Parse()
	targets := getopt.Args()
	if opts.IsDebug {
		opts.LogLevel = "debug"
	}

	if showVersion {
		fmt.Println("bogo", bogo.Version)
		return
	}
	if showHelp {
		fmt.Println("bogo", bogo.Version)
		getopt.Usage()
		return
	}

	c, cancel := bogo.NewDefaultContext(opts)

	c.Logger().Debug("targets:", targets)

	ch := make(chan interface{})

	keys := make([]string, 0, len(checks.Checkers))
	for k, x := range checks.Checkers {
		c.Logger().Debug("--- checker:", k, x)
		c.Logger().Info("starting checker ", k, "...")
		x.Run(c, ch)

		keys = append(keys, k)
	}
	c.Logger().Debug("checkers:", keys)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

main:
	for {
		select {
		case s := <-sig:
			c.Logger().Warnf("signal caught: %v", s)
			break main
		case m := <-ch:
			if pm, ok := m.(bogo.PingMessage); ok {
				c.Logger().Debug("ping message:", pm)
			} else {
				c.Logger().Debug("received:", m)
			}
		case <-time.After(50 * time.Millisecond):
		}
	}
	signal.Reset()

	c.Logger().Debug("cancelling the main context...")
	cancel()
	c.Logger().Debug("closing the channel...")
	close(ch)

	c.Logger().Debug("waiting for:", c.WG())
	c.WG().Wait()
}

/*
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
	run(&opts, targets)
}

func run(opts *bogo.Options, targets []string) {
	out := make(chan bogo.PingMessage)
	exporterLock := make(chan int)

	var exporter bogo.PingExproter
	switch opts.Exporter {
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
}
*/
