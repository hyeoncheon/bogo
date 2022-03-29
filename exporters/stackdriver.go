package exporters

import (
	"context"
	"fmt"

	"github.com/hyeoncheon/bogo"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type StackdriverExporter struct {
	instanceName string
	externalIP   string
	zone         string
}

var (
	avgRttMs = stats.Float64("ping_avgrtt", "average rtt in milliseconds", "ms")
	maxRttMs = stats.Float64("ping_maxrtt", "maximum rtt in milliseconds", "ms")
	lossRate = stats.Float64("ping_loss", "packet loss rate", "%")
)

func (e *StackdriverExporter) Initialize(in chan bogo.PingMessage, wait chan int) {
	bogo.Info("stackdriver exporter: initialize exporter...")
	c := bogo.NewMetadataClient()

	var err error
	e.instanceName, err = c.InstanceName()
	if err != nil {
		bogo.Err("could not get instance name: %v", err)
	}
	e.externalIP, err = c.ExternalIP()
	if err != nil {
		bogo.Err("could not get external IP: %v", err)
	}
	e.zone, err = c.Zone()
	if err != nil {
		bogo.Err("could not get zone: %v", err)
	}
	go e.run(in, wait)
}

func (e *StackdriverExporter) run(in chan bogo.PingMessage, wait chan int) {
	// register stackdriver view
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
		bogo.Err("could not register view: %v", err)
	}

	// register stackdriver view
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
		bogo.Err("could not register view: %v", err)
	}

	// create exporter instance for stackdriver
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		MetricPrefix: "custom.googleapis.com/bogo",
		GetMetricDisplayName: func(v *view.View) string {
			return fmt.Sprintf("bogo/%v", v.Name)
		},
	})
	if err != nil {
		bogo.Err("could not create exporter: %v", err)
	}
	defer exporter.Flush()

	if err := exporter.StartMetricsExporter(); err != nil {
		bogo.Err("could not start metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()

	defer bogo.Info("stackdriver: bye")

	ctx := context.Background()

	for {
		s, ok := <-in
		if !ok {
			wait <- 1
			return
		}

		bogo.Debug("stackdriver: got a input %v, %v\n", s, ok)
		if err := stats.RecordWithTags(ctx, []tag.Mutator{
			tag.Upsert(tag.MustNewKey("node"), e.instanceName),
			tag.Upsert(tag.MustNewKey("addr"), e.externalIP),
			tag.Upsert(tag.MustNewKey("zone"), e.zone),
			tag.Upsert(tag.MustNewKey("target"), s.Addr),
		},
			avgRttMs.M(float64(s.AvgRtt.Milliseconds())),
			lossRate.M(s.Loss),
		); err != nil {
			bogo.Err("stackdriver export: could not send ping stat")
		}
	}
}
