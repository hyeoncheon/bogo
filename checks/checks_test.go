package checks

import (
	"fmt"
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

func TestStartAll(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	opts.Checkers = []string{"heartbeat"}
	opts.CheckerOptions = map[string]common.PluginOptions{
		"heartbeat": map[string][]string{
			"interval": {"1"},
		},
	}
	ctx, _ := common.NewDefaultContext(&opts)
	defer ctx.Cancel()

	n := StartAll(ctx, &opts, ctx.Channel())
	r.Equal(1, n)

	out := <-ctx.Channel()
	r.NotNil(out)
	r.Equal("heartbeat", out.(string))
}

func TestStartAll_Error(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	opts.Checkers = []string{"heartbeat"}
	opts.CheckerOptions = map[string]common.PluginOptions{
		"heartbeat": map[string][]string{
			"interval": {"number"},
		},
	}
	ctx, _ := common.NewDefaultContext(&opts)
	defer ctx.Cancel()

	n := StartAll(ctx, &opts, ctx.Channel())
	r.Equal(0, n)

	r.Equal("--- &{{} [0 0 0]}", fmt.Sprintf("--- %v", ctx.WG()))
}
