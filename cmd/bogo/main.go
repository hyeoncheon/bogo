package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/checks"
	"github.com/hyeoncheon/bogo/internal/common"

	getopt "github.com/pborman/getopt/v2"
)

// main handles all options related tasks then calls run() with the options.
func main() {
	showVersion := false
	showHelp := false
	var copts string
	var eopts string

	opts := common.Options{
		IsDebug:   false,
		LogLevel:  "info",
		Checkers:  []string{},
		Exporters: []string{"stackdriver"},
	}
	getopt.FlagLong(&copts, "copts", 0, "checker options")
	getopt.FlagLong(&eopts, "eopts", 0, "exporter options")
	getopt.FlagLong(&opts.Checkers, "checker", 'c', "set checker")
	getopt.FlagLong(&opts.Exporters, "exporter", 'e', "set exporter")
	getopt.FlagLong(&opts.LogLevel, "log", 'l', "log level. (debug, info, warn, or error)")
	getopt.FlagLong(&opts.IsDebug, "debug", 'd', "debugging mode")
	getopt.FlagLong(&showVersion, "version", 'v', "show version")
	getopt.FlagLong(&showHelp, "help", 'h', "show help message")

	getopt.Parse()
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

	c, _ := common.NewDefaultContext(opts)
	logger := c.Logger()

	var err error
	opts.CheckerOptions, err = common.BuildPluginOptions(copts)
	if err != nil {
		logger.Fatal("could not build checker options:", err)
	}
	opts.ExporterOptions, err = common.BuildPluginOptions(eopts)
	if err != nil {
		logger.Fatal("could not build exporter options:", err)
	}

	logger.Debug("application opts:", opts)
	run(c, opts)
}

// run is the main thread
func run(c common.Context, opts common.Options) {
	logger := c.Logger()
	if c.Meta() == nil {
		logger.Warn("hey, it seems like I am on a legacy server or unsupported cloud!")
	}

	ch := make(chan interface{}) // communication channel for all plugins

	for k, x := range checks.Checkers {
		if len(opts.Checkers) > 0 && !common.Contains(opts.Checkers, k) {
			logger.Debugf("%v is not on the checker list. skipping...", k)
			continue
		}
		copts := opts.CheckerOptions[k]
		logger.Debug("--- checker:", k, x, copts)
		logger.Info("starting checker ", k, "...")
		x.Run(c, copts, ch)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

main:
	for {
		select {
		case s := <-sig:
			logger.Warnf("signal caught: %v", s)
			break main
		case m := <-ch:
			if pm, ok := m.(bogo.PingMessage); ok {
				logger.Debug("ping message:", pm)
			} else {
				logger.Debug("received:", m)
			}
		case <-time.After(50 * time.Millisecond):
		}
	}
	signal.Reset()

	logger.Debug("cancelling the main context...")
	c.Cancel()
	logger.Debug("closing the channel...")
	close(ch)
	logger.Debug("waiting for:", c.WG())
	c.WG().Wait()
}
