package meari

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hyeoncheon/bogo/internal/defaults"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

// TestNewDefaultServer tests the generator and if the singleton works fine.
func TestNewDefaultServer(t *testing.T) {
	r := require.New(t)

	s1 := NewDefaultServer(&Options{})
	r.NotNil(s1)

	s := NewDefaultServer(&Options{})
	r.NotNil(s)
	r.Equal(s1, s) // singleton by Once
	r.Equal(defaults.ServerAddress, s.Address())

	r.Nil(nil)
}

// TestNewDefaultServer_Address tests if Options.Address is working.
func TestNewDefaultServer_Address(t *testing.T) {
	r := require.New(t)

	serverOnce = sync.Once{}
	server = nil

	s := NewDefaultServer(&Options{Address: ":80"})
	r.NotNil(s)
	r.Equal(":80", s.Address())

	r.Nil(nil)
}

// TestDefaultServer_Functions tests a lifecycle of the server.
func TestDefaultServer_Functions(t *testing.T) {
	r := require.New(t)
	r.Nil(nil)

	serverOnce = sync.Once{}
	server = nil

	s := NewDefaultServer(&Options{})
	r.NotNil(s)
	r.Equal(defaults.ServerAddress, s.Address())

	go func() {
		err := s.Serve()
		r.Error(err)
		r.Contains(err.Error(), "Server closed")
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := "http://" + defaults.ServerAddress + "/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	r.NoError(err)
	resp, err := http.DefaultClient.Do(req)
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	r.NoError(s.Shutdown(ctx))
}
