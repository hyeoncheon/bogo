package exporters

import (
	"context"
	"fmt"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const (
	stackdriverMetricPrefix     = "custom.googleapis.com/bogo"
	stackdriverExporter         = "stackdriver"
	stackdriverExporterInterval = 1 * time.Minute
	recordTimeout               = 30 * time.Second
)

// RegisterStackdriver returns a new Exporter and it is used by StartAll().
func (*Exporter) RegisterStackdriver() *Exporter {
	return &Exporter{
		name:    stackdriverExporter,
		runFunc: stackdriverRunner,
	}
}

type reporter struct {
	instanceName string
	externalIP   string
	zone         string
}

var (
	avgRttMs = stats.Float64("ping_avgrtt", "average rtt in milliseconds", "ms")
	lossRate = stats.Float64("ping_loss", "packet loss rate", "%")
)

// stackdriverRunner is a runner function for the Stackdriver Exporter.
// It starts a go routine for the exporter and returns error status.
// The go routine runs forever until the context canceled, and will
// send received data to the Google Cloud Platform Cloud Monitoring
// (previously known as Stackdriver).
func stackdriverRunner(c common.Context, _ common.PluginOptions, in chan interface{}) error {
	logger := c.Logger().WithField("exporter", stackdriverExporter)

	r, err := getReporter(c)
	if err != nil {
		return err
	}

	if err := registerViews(); err != nil {
		return fmt.Errorf("could not register views: %w", err)
	}

	exporter, err := createAndStartExporter()
	if err != nil {
		return err
	}

	c.WG().Add(1)
	go func() { //nolint
		defer c.WG().Done()

		ticker := time.NewTicker(stackdriverExporterInterval)
		defer ticker.Stop()

		// defer for exporter
		defer exporter.Flush()
		defer exporter.StopMetricsExporter()

	infinite:
		for {
			select {
			case m, ok := <-in:
				if !ok {
					break infinite
				}

				if pm, ok := m.(bogo.PingMessage); ok {
					logger.Debugf("ping: %v", pm)
					if err := recordPingMessage(r, &pm); err != nil {
						logger.Errorf("message %v: %w", pm, err)
					}
				} else {
					logger.Warnf("unknown: %v", m)
				}
			case <-c.Done():
				break infinite
			}
		}
		logger.Infof("%s exporter exited", stackdriverExporter)
	}()

	return nil
}

// getReporter configures and returns a new instance of reporter if it runs
// on GCP. Currently, only GCE VM is supported for this feature.
func getReporter(c common.Context) (*reporter, error) {
	// currently, stackdriver exporter is only suppored on the GCE instance
	meta := c.Meta()
	if meta == nil || meta.WhereAmI() != common.GOOGLE {
		return nil, common.ErrNotOnGCE
	}

	return &reporter{
		instanceName: meta.InstanceName(),
		externalIP:   meta.ExternalIP(),
		zone:         meta.Zone(),
	}, nil
}

// registerViews configures and registers the views.
func registerViews() error {
	v := &view.View{
		Name:        "ping/rtt_average",
		Measure:     avgRttMs,
		Description: "ping average rtt",
		Aggregation: view.Distribution(0, 5, 10, 50, 100, 150, 200, 400), //nolint
		TagKeys: []tag.Key{
			tag.MustNewKey("node"),
			tag.MustNewKey("addr"),
			tag.MustNewKey("zone"),
			tag.MustNewKey("target"),
		},
	}
	if err := view.Register(v); err != nil {
		return err
	}

	vLoss := &view.View{
		Name:        "ping/packet_loss",
		Measure:     lossRate,
		Description: "ping packet loss rate",
		Aggregation: view.Distribution(0, 5, 10, 50, 100), //nolint
		TagKeys: []tag.Key{
			tag.MustNewKey("node"),
			tag.MustNewKey("addr"),
			tag.MustNewKey("zone"),
			tag.MustNewKey("target"),
		},
	}

	return view.Register(vLoss)
}

// createAndStartExporter starts the stackdriver exporter.
func createAndStartExporter() (*stackdriver.Exporter, error) {
	// create exporter instance for stackdriver
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		MetricPrefix: stackdriverMetricPrefix,
		GetMetricDisplayName: func(v *view.View) string {
			return fmt.Sprintf("bogo/%v", v.Name)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("could not create exporter: %w", err)
	}

	if err := exporter.StartMetricsExporter(); err != nil {
		return nil, fmt.Errorf("could not start metric exporter: %w", err)
	}

	return exporter, nil
}

// recordPingMessage sends the given ping message to the Cloud Monitoring.
func recordPingMessage(r *reporter, m *bogo.PingMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), recordTimeout)
	defer cancel()

	if err := stats.RecordWithTags(ctx,
		[]tag.Mutator{
			tag.Upsert(tag.MustNewKey("node"), r.instanceName),
			tag.Upsert(tag.MustNewKey("addr"), r.externalIP),
			tag.Upsert(tag.MustNewKey("zone"), r.zone),
			tag.Upsert(tag.MustNewKey("target"), m.Addr),
		},
		avgRttMs.M(float64(m.AvgRtt.Milliseconds())),
		lossRate.M(m.Loss),
	); err != nil {
		return fmt.Errorf("could not send ping stat: %w", err)
	}

	return nil
}
