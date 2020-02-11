package exporters

import (
	"context"
	"fmt"
	"net/http"
	"prober/checks"

	"cloud.google.com/go/compute/metadata"
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
)

func (e *StackdriverExporter) Initialize(in chan checks.PingMessage, wait chan int) {
	fmt.Printf("stackdriver exporter: initialize exporter...\n")
	c := metadata.NewClient(&http.Client{Transport: userAgentTransport{
		userAgent: "prober-exporter-stackdriver-0.0.1",
		base:      http.DefaultTransport,
	}})
	var err error
	e.instanceName, err = c.InstanceName()
	if err != nil {
		fmt.Printf("could not get instance name: %v", err)
	}
	e.externalIP, err = c.ExternalIP()
	if err != nil {
		fmt.Printf("could not get external IP: %v", err)
	}
	e.zone, err = c.Zone()
	if err != nil {
		fmt.Printf("could not get zone: %v", err)
	}
	go e.run(in, wait)
}

func (e *StackdriverExporter) run(in chan checks.PingMessage, wait chan int) {
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
		fmt.Println("could not register view:", err)
	}

	// create exporter instance for stackdriver
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		MetricPrefix: "custom.googleapis.com/prober",
		GetMetricDisplayName: func(v *view.View) string {
			return fmt.Sprintf("prober/%v", v.Name)
		},
	})
	if err != nil {
		fmt.Println("could not create exporter:", err)
	}
	defer exporter.Flush()

	if err := exporter.StartMetricsExporter(); err != nil {
		fmt.Println("could not start metric exporter:", err)
	}
	defer exporter.StopMetricsExporter()

	defer fmt.Printf("stackdriver: bye\n")

	ctx := context.Background()

	for {
		s, ok := <-in
		if !ok {
			wait <- 1
			return
		}

		fmt.Printf("stackdriver: got a input %v, %v\n", s, ok)
		stats.RecordWithTags(ctx, []tag.Mutator{
			tag.Upsert(tag.MustNewKey("node"), e.instanceName),
			tag.Upsert(tag.MustNewKey("addr"), e.externalIP),
			tag.Upsert(tag.MustNewKey("zone"), e.zone),
			tag.Upsert(tag.MustNewKey("target"), s.Addr),
		}, avgRttMs.M(float64(s.AvgRtt.Milliseconds())))
	}
}

// from https://github.com/googleapis/google-cloud-go/blob/master/compute/metadata/examples_test.go
// userAgentTransport sets the User-Agent header before calling base.
type userAgentTransport struct {
	userAgent string
	base      http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}
