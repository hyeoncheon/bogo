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

func (*Exporter) RegisterStdout() *Exporter {
	return &Exporter{
		name:    stdoutExporter,
		runFunc: stdoutRunner,
	}
}

func stdoutRunner(c common.Context, _ common.PluginOptions, in chan interface{}) error {
	logger := c.Logger().WithField("exporter", stdoutExporter)
	c.WG().Add(1)
	go func() {
		defer c.WG().Done()

		ticker := time.NewTicker(stdoutExporterInterval)
		defer ticker.Stop()

	infinite:
		for {
			m, ok := <-in
			if !ok {
				break infinite
			}

			if pm, ok := m.(bogo.PingMessage); ok {
				logger.Infof("ping: %v", pm)
			} else {
				logger.Warnf("unknown: %v", m)
			}
		}
		logger.Infof("%s exporter exited", stdoutExporter)
	}()
	return nil
}
