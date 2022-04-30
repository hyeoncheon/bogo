package exporters

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"

	"github.com/stretchr/testify/require"
)

func TestStartAll(t *testing.T) {
	r := require.New(t)
	opts := common.DefaultOptions()
	opts.Exporters = []string{"stdout"}
	opts.ExporterOptions = map[string]common.PluginOptions{
		"stdout": map[string][]string{
			"interval": {"1"},
		},
	}
	ctx, _ := common.NewDefaultContext(&opts)

	n := StartAll(ctx, &opts, ctx.Channel())
	r.Equal(1, n)

	ctx.Channel() <- "message"

	ctx.Cancel()
}

func TestStartAll_MoreThanOne(t *testing.T) {
	r := require.New(t)
	opts := common.DefaultOptions()
	opts.Exporters = []string{"stdout", "stackdriver"}
	opts.ExporterOptions = map[string]common.PluginOptions{}
	ctx, _ := common.NewDefaultContext(&opts)

	n := StartAll(ctx, &opts, ctx.Channel())
	r.Equal(0, n)

	ctx.Cancel()
}
