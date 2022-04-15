package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/checks"
	"github.com/hyeoncheon/bogo/exporters"
	"github.com/hyeoncheon/bogo/handlers"
	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/hyeoncheon/bogo/meari"

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
		Address:   "127.0.0.1:6090",
	}
	getopt.FlagLong(&copts, "copts", 0, "checker options")
	getopt.FlagLong(&eopts, "eopts", 0, "exporter options")
	getopt.FlagLong(&opts.Checkers, "checker", 'c', "set checker")
	getopt.FlagLong(&opts.Exporters, "exporter", 'e', "set exporter")
	getopt.FlagLong(&opts.Address, "address", 'a', "webserver's listen address")
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

	c, _ := common.NewDefaultContext(&opts)
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
	run(c, &opts)
}

// run is the main thread
func run(c common.Context, opts *common.Options) {
	logger := c.Logger()
	if c.Meta() == nil {
		logger.Warn("hey, it seems like I am on a legacy server or unsupported cloud!")
	}

	ch := make(chan interface{}) // communication channel for all plugins
	defer func() {
		if r := recover(); r != nil {
			logger.Fatalf("panic: %v", r)
		}
	}()

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

	for k, x := range exporters.Exporters {
		if len(opts.Exporters) > 0 && !common.Contains(opts.Exporters, k) {
			logger.Debugf("%v is not on the exporter list. skipping...", k)
			continue
		}
		copts := opts.ExporterOptions[k]
		logger.Debug("--- exporter:", k, x, copts)
		logger.Info("starting exporter ", k, "...")
		x.Run(c, copts, ch)
	}

	serverOpts := &meari.Options{
		Logger:  logger.WithField("component", "web"),
		Address: opts.Address,
	}
	server := meari.New(serverOpts)
	if server != nil {
		for p, handler := range handlers.AllHanders() {
			switch handler.Method {
			case http.MethodGet:
				logger.Debugf("register handler for 'GET %v'...", p)
				server.GET(p, handler.Handler)
			default:
				logger.Errorf("unsupported method for %v", handler)
			}
		}

		c.WG().Add(1)
		go func() {
			defer c.WG().Done()

			err := server.Start()

			if err == http.ErrServerClosed {
				logger.Info("webserver closed successfully")
			} else {
				logger.Error("unexpected error: ", err)
			}
		}()
	} else {
		logger.Error("could not initiate the web server: ", serverOpts)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

main:
	for {
		select {
		case s := <-sig:
			logger.Warnf("signal caught: %v", s)
			break main
		case <-time.After(500 * time.Millisecond):
		}
	}
	signal.Reset()

	logger.Info("shutting down webserver...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("could not gracefully shutdown the web server: ", err)
	}

	logger.Debug("cancelling the main context...")
	c.Cancel()
	logger.Debug("closing the channel...")
	close(ch)
	logger.Debug("waiting for routines: ", c.WG())
	c.WG().Wait()
}
