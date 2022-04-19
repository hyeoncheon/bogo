package checks

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

func TestRegisterHeartbeat(t *testing.T) {
	r := require.New(t)

	p := (&Checker{}).RegisterHeartbeat()
	r.IsType(&Checker{}, p)
	r.Implements((*common.Plugin)(nil), p)

	r.Equal(heartbeatChecker, p.Name())
}

func TestHeartbeatRunner(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, cancel := common.NewDefaultContext(&opts)
	ch := make(chan interface{})

	o := common.PluginOptions{
		"interval": []string{"1"},
	}

	r.NoError(heartbeatRunner(c, o, ch))
	get := <-ch
	r.NotNil(get)
	r.Equal("heartbeat", get.(string))
	cancel()

	c.WG().Wait()
}

func TestHeartbeatRunner_Panic(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, cancel := common.NewDefaultContext(&opts)
	defer cancel()
	ch := make(chan interface{})

	o := common.PluginOptions{
		"interval": []string{"1"},
	}

	r.NoError(heartbeatRunner(c, o, ch))
	close(ch)

	c.WG().Wait()
}

func TestHeartbeatRunner_WrongOption(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	ch := make(chan interface{})

	o := common.PluginOptions{
		"interval": []string{"number"},
	}

	err := heartbeatRunner(c, o, ch)
	r.ErrorContains(err, "invalid option value")
}
