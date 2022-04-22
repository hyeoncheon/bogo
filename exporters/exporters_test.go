package exporters

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
)

func TestStartAll(t *testing.T) {
	opts := common.DefaultOptions()
	opts.Exporters = []string{"stdout"}
	opts.ExporterOptions = map[string]common.PluginOptions{
		"stdout": map[string][]string{
			"interval": {"1"},
		},
	}
	ctx, _ := common.NewDefaultContext(&opts)

	StartAll(ctx, &opts, ctx.Channel())

	ctx.Channel() <- "message"

	ctx.Cancel()
}
