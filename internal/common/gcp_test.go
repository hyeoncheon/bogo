package common

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"cloud.google.com/go/compute/metadata"
	"github.com/stretchr/testify/require"
)

func TestNewGCEMetaClient_SimGCE(t *testing.T) {
	r := require.New(t)

	// mocking with true
	originalFunc := gceMetaOnGCE
	oGceClientOnce := gceClientOnce
	oGceClient := gceClient
	defer func() {
		gceMetaOnGCE = originalFunc
		gceClientOnce = oGceClientOnce
		gceClient = oGceClient
	}()
	gceClientOnce = sync.Once{}
	gceClient = nil
	gceMetaOnGCE = func() bool {
		return true
	}

	m := NewGCEMetaClient(&defaultContext{})
	r.IsType((*GCEClient)(nil), m)
	r.Implements((*MetaClient)(nil), m)

	res := m.WhereAmI()
	r.Equal(GOOGLE, res)
}

func TestNewGCEMetaClient_NoGCE(t *testing.T) {
	r := require.New(t)

	// mocking with false
	originalFunc := gceMetaOnGCE
	oGceClientOnce := gceClientOnce
	oGceClient := gceClient
	defer func() {
		gceMetaOnGCE = originalFunc
		gceClientOnce = oGceClientOnce
		gceClient = oGceClient
	}()
	gceClientOnce = sync.Once{}
	gceClient = nil
	gceMetaOnGCE = func() bool {
		return false
	}

	m := NewGCEMetaClient(&defaultContext{})
	r.Nil(m)

	// actually this could be separated test but just keep it here
	m = &GCEClient{Client: nil, logger: nil}
	res := m.WhereAmI()
	r.Equal(NOWHERE, res)
}

func TestGCEClient_InstanceName(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaInstanceName
	defer func() {
		gceMetaInstanceName = originalFunc
	}()

	gceMetaInstanceName = func(*metadata.Client) (string, error) {
		return "instance-1", nil
	}
	res := meta.InstanceName()
	r.Equal("instance-1", res)

	gceMetaInstanceName = func(*metadata.Client) (string, error) {
		return "instance-1", errors.New("error")
	}
	res = meta.InstanceName()
	r.Equal(UNKNOWN, res)
}

func TestGCEClient_ExternalIP(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaExternalIP
	defer func() {
		gceMetaExternalIP = originalFunc
	}()

	gceMetaExternalIP = func(*metadata.Client) (string, error) {
		return "203.0.113.1", nil
	}
	res := meta.ExternalIP()
	r.Equal("203.0.113.1", res)

	gceMetaExternalIP = func(*metadata.Client) (string, error) {
		return "203.0.113.1", errors.New("error")
	}
	res = meta.ExternalIP()
	r.Equal(UNKNOWN, res)
}

func TestGCEClient_Zone(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaZone
	defer func() {
		gceMetaZone = originalFunc
	}()

	gceMetaZone = func(*metadata.Client) (string, error) {
		return "asia-northeast3-z", nil
	}
	res := meta.Zone()
	r.Equal("asia-northeast3-z", res)

	gceMetaZone = func(*metadata.Client) (string, error) {
		return "asia-northeast3-z", errors.New("error")
	}
	res = meta.Zone()
	r.Equal(UNKNOWN, res)
}

func TestGCEClient_AttributeValue(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaInstanceAttributeValue
	originalFunc2 := gceMetaProjectAttributeValue
	defer func() {
		gceMetaInstanceAttributeValue = originalFunc
		gceMetaProjectAttributeValue = originalFunc2
	}()

	// instance has the meta
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "instance's meta value", nil
	}
	gceMetaProjectAttributeValue = func(*metadata.Client, string) (string, error) {
		return "project's meta value", nil
	}
	res := meta.AttributeValue("key")
	r.Equal("instance's meta value", res)

	// instance has no meta defined but the project have one
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "", errors.New("nodata")
	}
	gceMetaProjectAttributeValue = func(*metadata.Client, string) (string, error) {
		return "project's meta value", nil
	}
	res = meta.AttributeValue("key")
	r.Equal("project's meta value", res)

	// instance has the meta defined but the value is empty.
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "", nil
	}
	gceMetaProjectAttributeValue = func(*metadata.Client, string) (string, error) {
		return "project's meta value", nil
	}
	res = meta.AttributeValue("key")
	r.Equal("", res)

	// they both have no such metadata
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "", errors.New("nodata")
	}
	gceMetaProjectAttributeValue = func(*metadata.Client, string) (string, error) {
		return "", errors.New("nodata")
	}
	res = meta.AttributeValue("key")
	r.Equal("", res)

	// returns empty string with nil
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "", nil
	}
	res = meta.AttributeValue("key")
	r.Equal("", res)
}

func TestGCEClient_AttributeCSV(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaInstanceAttributeValue
	defer func() {
		gceMetaInstanceAttributeValue = originalFunc
	}()

	// instance has the meta
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "hey, bulldog oh,yeah", nil
	}
	res := meta.AttributeCSV("key")
	r.EqualValues([]string{"hey", "bulldog oh", "yeah"}, res)
}

func TestGCEClient_AttributeSSV(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaInstanceAttributeValue
	defer func() {
		gceMetaInstanceAttributeValue = originalFunc
	}()

	// instance has the meta
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "hey, bulldog oh,yeah", nil
	}
	res := meta.AttributeSSV("key")
	r.EqualValues([]string{"hey,", "bulldog", "oh,yeah"}, res)
}

func TestGCEClient_AttributeValues(t *testing.T) {
	r := require.New(t)

	meta := &GCEClient{Client: &metadata.Client{}, logger: nil}

	// mocking with false
	originalFunc := gceMetaInstanceAttributeValue
	defer func() {
		gceMetaInstanceAttributeValue = originalFunc
	}()

	// instance has the meta
	gceMetaInstanceAttributeValue = func(*metadata.Client, string) (string, error) {
		return "hey, bulldog oh,yeah", nil
	}
	res := meta.AttributeValues("key")
	r.EqualValues([]string{"hey", "bulldog", "oh", "yeah"}, res)
}

func TestGCEClient_RoundTrip(t *testing.T) {
	r := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ua := req.Header.Get("User-Agent")
		rw.Write([]byte(ua))
	}))
	defer server.Close()

	client := server.Client()
	client.Transport = userAgentTransport{
		userAgent: "my client/1.0",
		base:      http.DefaultTransport,
	}

	resp, err := client.Get(server.URL)
	r.NoError(err)
	r.NotNil(resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	r.NoError(err)
	r.Equal([]byte("my client/1.0"), body)
}
