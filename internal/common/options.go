package common

import "github.com/hyeoncheon/bogo/internal/defaults"

// Options is a struct to store command line options
type Options struct {
	IsDebug  bool
	LogLevel string
	// Checkers constains selected checkers. All available checkers will be
	// selected if this field is empty.
	Checkers []string
	// Exporters contains selected exporters. All available exporters will be
	// selected if this field is empty.
	Exporters []string

	// Address is a listen address for embedded webserver
	Address string

	CheckerOptions  map[string]PluginOptions
	ExporterOptions map[string]PluginOptions
}

func DefaultOptions() Options {
	return Options{
		IsDebug:   false,
		LogLevel:  "info",
		Checkers:  []string{},
		Exporters: []string{"stackdriver"},
		Address:   defaults.ServerAddress,
	}
}
