package common

import (
	"net/http"
	"strings"
	"sync"

	"github.com/hyeoncheon/bogo"

	"cloud.google.com/go/compute/metadata"
)

// asset GCEClient for MetaClient iplemetations
var _ MetaClient = &GCEClient{}

// GCEClient is a struct for handling GCE metadata client
type GCEClient struct {
	*metadata.Client
	logger Logger
}

var (
	gceClientOnce sync.Once
	gceClient     MetaClient
)

// NewGCEMetaClient tests if the application is running on a GCE instance,
// then returns the `GCEClient` as `MetaClient`.
func NewGCEMetaClient(c Context) MetaClient {
	gceClientOnce.Do(func() {
		if metadata.OnGCE() {
			gceClient = &GCEClient{
				Client: newGoogleCloudMetadataClient(),
				logger: c.Logger().WithField("meta", "gcp"),
			}
		}
	})
	return gceClient
}

func (m *GCEClient) WhereAmI() string {
	return "Google"
}

func (m *GCEClient) AttributeValue(key string) string {
	value, err := m.InstanceAttributeValue(key)
	if err != nil || len(value) == 0 {
		m.logger.Debugf("no '%v' in instance attributes.", key)
		value, err = m.ProjectAttributeValue(key)
		if err != nil || len(value) == 0 {
			m.logger.Debugf("no '%v' in project attributes.", key)
		}
	}
	return value
}

func (m *GCEClient) AttributeCSV(s string) []string {
	result := []string{}
	for _, t := range strings.Split(m.AttributeValue(s), ",") {
		result = append(result, strings.TrimSpace(t))
	}
	return result
}

func (m *GCEClient) AttributeSSV(s string) []string {
	result := []string{}
	for _, t := range strings.Split(m.AttributeValue(s), " ") {
		result = append(result, strings.TrimSpace(t))
	}
	return result
}

func (m *GCEClient) AttributeValues(s string) []string {
	result := []string{}
	for _, t := range strings.Split(m.AttributeValue(s), ",") {
		for _, t := range strings.Split(strings.TrimSpace(t), " ") {
			result = append(result, strings.TrimSpace(t))
		}
	}
	return result
}

// from https://github.com/googleapis/google-cloud-go/blob/master/compute/metadata/examples_test.go

// newGoogleCloudMetadataClient
func newGoogleCloudMetadataClient() *metadata.Client {
	m := metadata.NewClient(&http.Client{
		Transport: userAgentTransport{
			userAgent: "bogo/" + bogo.Version,
			base:      http.DefaultTransport,
		},
	})
	return m
}

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
