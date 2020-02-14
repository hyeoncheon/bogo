package prober

import (
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
)

func NewMetadataClient() *metadata.Client {
	return metadata.NewClient(&http.Client{Transport: userAgentTransport{
		userAgent: "prober-0.0.2",
		base:      http.DefaultTransport,
	}})
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

func Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Warn(format string, args ...interface{}) {
	Err(format+"\n", args...)
}

func Err(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func Debug(format string, args ...interface{}) {
}
