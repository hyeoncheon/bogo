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

	opts := common.DefaultOptions()

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

	defer func() {
		if r := recover(); r != nil {
			logger.Fatalf("panic: %v", r)
		}
	}()

	checks.StartAll(c, opts, c.Channel())
	exporters.StartAll(c, opts, c.Channel())

	server, err := startWebRoutine(c, opts)
	if err != nil {
		logger.Errorf("could not start the webserver: ", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	s := <-sig
	logger.Warnf("signal caught: %v", s)
	signal.Reset()

	logger.Info("shutting down webserver...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("could not gracefully shutdown the web server: ", err)
	}

	c.Cancel()
}

func startWebRoutine(c common.Context, opts *common.Options) (meari.Server, error) {
	logger := c.Logger()

	server, err := meari.NewServer(c, opts)
	if err != nil {
		return nil, err
	}

	c.WG().Add(1)
	go func() {
		defer c.WG().Done()

		err := server.Start()
		if err == http.ErrServerClosed {
			logger.Info("webserver closed")
		} else {
			logger.Error("unexpected error: ", err)
		}
	}()
	return server, nil
}
