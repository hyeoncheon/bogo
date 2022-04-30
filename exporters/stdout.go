package exporters

import (
	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"
)

const (
	stdoutExporter = "stdout"
)

// RegisterStdout returns a new Exporter and it is used by StartAll().
func (*Exporter) RegisterStdout() *Exporter {
	return &Exporter{
		name:    stdoutExporter,
		runFunc: stdoutRunner,
	}
}

// stdoutRunner is a runner function for the Stdout Exporter. It starts a go
// routine for the exporter and returns the status, then the exporter runs
// forever and will print out catched message until the context is canceled.
func stdoutRunner(c common.Context, _ common.PluginOptions, in chan interface{}) error {
	logger := c.Logger().WithField("exporter", stdoutExporter)

	c.WG().Add(1)
	go func() { //nolint
		defer c.WG().Done()

	infinite:
		for {
			select {
			case m, ok := <-in:
				if !ok {
					break infinite
				}

				if pm, ok := m.(bogo.PingMessage); ok {
					logger.Infof("ping: %v", pm)
				} else {
					logger.Warnf("unknown: %v", m)
				}
			case <-c.Done():
				break infinite
			}
		}
		logger.Infof("%s exporter exited", stdoutExporter)
	}()

	return nil
}
