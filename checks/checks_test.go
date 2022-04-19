package checks

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

func TestStartAll(t *testing.T) {
	r := require.New(t)
	r.Nil(nil)

	opts := common.DefaultOptions()
	opts.Checkers = []string{"heartbeat"}
	opts.CheckerOptions = map[string]common.PluginOptions{
		"heartbeat": map[string][]string{
			"interval": {"1"},
		},
	}
	ctx, _ := common.NewDefaultContext(&opts)

	ch := make(chan interface{})
	defer close(ch)

	StartAll(ctx, &opts, ch)

	out := <-ch
	r.NotNil(out)
	r.Equal("heartbeat", out.(string))

	ctx.Cancel()
	ctx.WG().Wait()
}

func TestStartAll_Error(t *testing.T) {
	r := require.New(t)
	r.Nil(nil)

	opts := common.DefaultOptions()
	opts.Checkers = []string{"heartbeat"}
	opts.CheckerOptions = map[string]common.PluginOptions{
		"heartbeat": map[string][]string{
			"interval": {"number"},
		},
	}
	ctx, _ := common.NewDefaultContext(&opts)

	ch := make(chan interface{})
	defer close(ch)

	StartAll(ctx, &opts, ch)

	ctx.Cancel()
	ctx.WG().Wait()
}
