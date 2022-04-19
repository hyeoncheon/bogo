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

func TestNewDefaultServer_Address(t *testing.T) {
	r := require.New(t)

	serverOnceOrig := serverOnce
	serverOrig := server
	defer func() {
		serverOnce = serverOnceOrig
		server = serverOrig
	}()
	serverOnce = sync.Once{}
	server = nil

	s := NewDefaultServer(&Options{Address: ":80"})
	r.NotNil(s)
	r.Equal(":80", s.Address())

	r.Nil(nil)
}

func TestDefaultServer_Functions(t *testing.T) {
	r := require.New(t)
	r.Nil(nil)

	s := NewDefaultServer(&Options{})
	r.NotNil(s)
	r.Equal(defaults.ServerAddress, s.Address())

	go func() {
		err := s.Start()
		r.Error(err)
		r.Contains(err.Error(), "Server closed")
	}()

	time.Sleep(100 * time.Millisecond)

	c := http.Client{}
	resp, err := c.Get("http://" + defaults.ServerAddress + "/")
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	r.NoError(s.Shutdown(ctx))
}
