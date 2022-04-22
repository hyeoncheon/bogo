package exporters

import (
	"testing"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

func TestRegisterStdout(t *testing.T) {
	r := require.New(t)

	p := (&Exporter{}).RegisterStdout()
	r.IsType(&Exporter{}, p)
	r.Implements((*common.Plugin)(nil), p)

	r.Equal(stdoutExporter, p.Name())
}

func TestStdoutRunner(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{}

	r.NoError(stdoutRunner(c, o, c.Channel()))
	time.Sleep(100 * time.Millisecond)
	c.Channel() <- "test"
	c.Channel() <- bogo.PingMessage{}
	time.Sleep(100 * time.Millisecond)

	c.Cancel()
}

func TestStdoutRunner_Closed(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{}

	r.NoError(stdoutRunner(c, o, c.Channel()))
	time.Sleep(100 * time.Millisecond)
	close(c.Channel())
	time.Sleep(100 * time.Millisecond)

	c.Cancel() // will cause a panic, defer recovery added for this
}
