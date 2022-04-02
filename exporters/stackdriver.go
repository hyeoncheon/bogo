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
	stackdriverExporter         = "stackdriver"
	stackdriverExporterInterval = 1 * time.Minute
)

func (x *Exporter) Stackdriver() error {
	x.Name = stackdriverExporter
	x.Run = stackdriverRunner
	return nil
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

func stackdriverRunner(c common.Context, opts common.PluginOptions, in chan interface{}) error {
	logger := c.Logger().WithField("exporter", stackdriverExporter)

	meta := c.Meta()
	if meta == nil || meta.WhereAmI() != "Google" {
		logger.Error("could not get Google Cloud meta client!")
		return fmt.Errorf("count not get Google Cloud metadata client")
		// TODO: need to implement remote reporting
	}

	r := &reporter{
		instanceName: meta.InstanceName(),
		externalIP:   meta.ExternalIP(),
		zone:         meta.Zone(),
	}

	c.WG().Add(1)
	go func() {
		defer c.WG().Done()
		ticker := time.NewTicker(stackdriverExporterInterval)
		defer ticker.Stop()

		if err := registerViews(); err != nil {
			logger.Error("could not register views:", err)
			return
		}

		// create exporter instance for stackdriver
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			MetricPrefix: "custom.googleapis.com/bogo",
			GetMetricDisplayName: func(v *view.View) string {
				return fmt.Sprintf("bogo/%v", v.Name)
			},
		})
		if err != nil {
			logger.Errorf("could not create exporter: %v", err)
			return
		}
		defer exporter.Flush()

		if err := exporter.StartMetricsExporter(); err != nil {
			logger.Errorf("could not start metric exporter: %v", err)
			return
		}
		defer exporter.StopMetricsExporter()

		ctx := context.Background()

	infinit:
		for {
			rm, ok := <-in
			if !ok {
				break infinit
			}
			if pm, ok := rm.(bogo.PingMessage); ok {
				logger.Debugf("ping: %v", pm)
				if err := stats.RecordWithTags(ctx, []tag.Mutator{
					tag.Upsert(tag.MustNewKey("node"), r.instanceName),
					tag.Upsert(tag.MustNewKey("addr"), r.externalIP),
					tag.Upsert(tag.MustNewKey("zone"), r.zone),
					tag.Upsert(tag.MustNewKey("target"), pm.Addr),
				},
					avgRttMs.M(float64(pm.AvgRtt.Milliseconds())),
					lossRate.M(pm.Loss),
				); err != nil {
					logger.Error("could not send ping stat:", err)
				}
			} else {
				logger.Infof("known: %v (%v)", rm, ok)
			}
		}
		logger.Info(stackdriverExporter, " done.")
	}()
	return nil
}

func registerViews() error {
	v := &view.View{
		Name:        "ping/rtt_average",
		Measure:     avgRttMs,
		Description: "ping average rtt",
		Aggregation: view.Distribution(0, 5, 10, 50, 100, 150, 200, 400),
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
		Aggregation: view.Distribution(0, 5, 10, 50, 100),
		TagKeys: []tag.Key{
			tag.MustNewKey("node"),
			tag.MustNewKey("addr"),
			tag.MustNewKey("zone"),
			tag.MustNewKey("target"),
		},
	}
	if err := view.Register(vLoss); err != nil {
		return err
	}
	return nil
}
