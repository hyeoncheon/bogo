package meari

import (
	"sync"
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

// TestNewServer tests NewServer with default Options and default Context.
func TestNewServer(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	ctx, _ := common.NewDefaultContext(&opts)

	s, err := NewServer(ctx, &opts)
	r.NoError(err)
	r.NotNil(s)
}

func TestNewServer_Address(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	ctx, _ := common.NewDefaultContext(&opts)

	serverOnce = sync.Once{}

	opts.Address = "address"
	s, err := NewServer(ctx, &opts)
	r.NoError(err)
	r.NotNil(s)
	r.EqualValues("address", s.Address())
}

func TestNewServer_Once(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	ctx, _ := common.NewDefaultContext(&opts)

	serverOnce = sync.Once{}

	s1, err := NewServer(ctx, &opts)
	r.NoError(err)
	r.NotNil(s1)

	opts.Address = "address"
	s2, err := NewServer(ctx, &opts)
	r.NoError(err)
	r.NotNil(s2)

	r.EqualValues(s1, s2)                     // singleton
	r.EqualValues(s1.Address(), s2.Address()) // singleton
}
