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
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{
		"interval": []string{"1"},
	}

	r.NoError(heartbeatRunner(c, o, c.Channel()))
	get := <-c.Channel()
	r.NotNil(get)
	r.Equal("heartbeat", get.(string))

	c.Cancel()
}

/* changed the shutdown handling gracefully
func TestHeartbeatRunner_Panic(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{
		"interval": []string{"1"},
	}

	r.NoError(heartbeatRunner(c, o, c.Channel()))
	close(c.Channel())

	c.Cancel()
}
*/

func TestHeartbeatRunner_WrongOption(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	defer c.Cancel()

	o := common.PluginOptions{
		"interval": []string{"number"},
	}

	err := heartbeatRunner(c, o, c.Channel())
	r.ErrorContains(err, "invalid option value")
}
