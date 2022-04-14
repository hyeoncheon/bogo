package exporters

import (
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	stdoutExporter         = "stdout"
	stdoutExporterInterval = 1 * time.Minute
)

func (x *Exporter) Stdout() error {
	x.Name = stdoutExporter
	x.Run = stdoutRunner
	return nil
}

func stdoutRunner(c common.Context, opts common.PluginOptions, in chan interface{}) error {
	logger := c.Logger().WithField("exporter", stdoutExporter)
	c.WG().Add(1)
	go func() {
		defer c.WG().Done()
		ticker := time.NewTicker(stdoutExporterInterval)
		defer ticker.Stop()

	infinite:
		for {
			m, ok := <-in
			if pm, ok := m.(bogo.PingMessage); ok {
				logger.Infof("ping: %v", pm)
			} else {
				logger.Infof("unknown: %v", m)
			}
			if !ok {
				break infinite
			}
		}
		logger.Info(stdoutExporter, " done.")
	}()
	return nil
}
