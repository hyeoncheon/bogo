package common

import (
	"net/http"
	"strings"
	"sync"

	"github.com/hyeoncheon/bogo"

	"cloud.google.com/go/compute/metadata"
)

const GOOGLE = "Google"

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

// gceMetaOnGCE is a function pointer to metadata.OnGCE().
// Use it to make unit test easier.
var gceMetaOnGCE = metadata.OnGCE

// NewGCEMetaClient tests if the application is running on a GCE instance,
// then returns the `GCEClient` as `MetaClient`.
func NewGCEMetaClient(c Context) MetaClient {
	logger := c.Logger()
	if logger == nil {
		logger = NewDefaultLogger("info")
	}

	gceClientOnce.Do(func() {
		if gceMetaOnGCE() {
			gceClient = &GCEClient{
				Client: newGoogleCloudMetadataClient(),
				logger: logger.WithField("meta", "gcp"),
			}
		}
	})
	return gceClient
}

// WhereAmI returns the name of CSP.
func (m *GCEClient) WhereAmI() string {
	if m.Client != nil {
		return GOOGLE
	}
	return NOWHERE
}

// gceMetaInstanceName is a function pointer to metadata.InstanceName().
// Use it to make unit test easier.
var gceMetaInstanceName = (*metadata.Client).InstanceName

// InstanceName returns the current VM's instance name string.
func (m *GCEClient) InstanceName() string {
	ret, err := gceMetaInstanceName(m.Client)
	if err != nil {
		return UNKNOWN
	}
	return ret
}

// gceMetaExternalIP is a function pointer to metadata.ExternalIP().
// Use it to make unit test easier.
var gceMetaExternalIP = (*metadata.Client).ExternalIP

// ExternalIP returns the instance's primary external (public) IP address.
func (m *GCEClient) ExternalIP() string {
	ret, err := gceMetaExternalIP(m.Client)
	if err != nil {
		return UNKNOWN
	}
	return ret
}

// gceMetaZone is a function pointer to metadata.Zone().
// Use it to make unit test easier.
var gceMetaZone = (*metadata.Client).Zone

// Zone returns the current VM's zone, such as "asia-northeast3-a".
func (m *GCEClient) Zone() string {
	ret, err := gceMetaZone(m.Client)
	if err != nil {
		return UNKNOWN
	}
	return ret
}

// gceMetaInstanceAttributeValue is a function pointer to metadata.InstanceAttributeValue().
// Use it to make unit test easier.
var gceMetaInstanceAttributeValue = (*metadata.Client).InstanceAttributeValue

// gceMetaProjectAttributeValue is a function pointer to metadata.ProjectAttributeValue().
// Use it to make unit test easier.
var gceMetaProjectAttributeValue = (*metadata.Client).ProjectAttributeValue

// AttributeValue returns the raw metadata stored for the instance. It returns empty
// string if the metadata with the key exists but the value is empty. When the instance
// has no the metadata defined, it will returns the project's metadata.
func (m *GCEClient) AttributeValue(key string) string {
	if m.logger == nil {
		m.logger = NewDefaultLogger("info").WithField("meta", "gcp")
	}

	value, err := gceMetaInstanceAttributeValue(m.Client, key)
	if err != nil {
		m.logger.Debugf("no '%v' in instance attributes.", key)
		value, err = gceMetaProjectAttributeValue(m.Client, key)
		if err != nil {
			m.logger.Debugf("no '%v' in project attributes.", key)
		}
	}
	return value
}

// AttributeCSV returns the metadata as an array of strings. The raw value
// will be treated as comma separated values. So if the raw metadata is "oh,
// little darling", the returned array will be consist of "oh" and "little
// darling".
func (m *GCEClient) AttributeCSV(s string) []string {
	result := []string{}
	for _, t := range strings.Split(m.AttributeValue(s), ",") {
		result = append(result, strings.TrimSpace(t))
	}
	return result
}

// AttributeSSV returns the metadata as an array of strings. The raw value
// will be treated as space separated values. So if the raw metadata is "oh,
// little darling", the returned array will be consist of "oh,", "little",
// and "darling".
func (m *GCEClient) AttributeSSV(s string) []string {
	result := []string{}
	for _, t := range strings.Split(m.AttributeValue(s), " ") {
		result = append(result, strings.TrimSpace(t))
	}
	return result
}

// AttributeValues returns the metadata as an array of strings. The raw value
// will be treated as both comma and space separated values. So if the raw
// metadata is "oh, little darling", the returned array will be consist of
// "oh", "little", and "darling".
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
			userAgent: bogo.Name + "/" + bogo.Version,
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
