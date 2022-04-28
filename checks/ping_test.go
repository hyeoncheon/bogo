package checks

import (
	"testing"
	"time"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

func TestRegisterPing(t *testing.T) {
	r := require.New(t)

	p := (&Checker{}).RegisterPing()
	r.IsType(&Checker{}, p)
	r.Implements((*common.Plugin)(nil), p)

	r.Equal(pingChecker, p.Name())
}

func TestPingRunner(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{
		"targets":        []string{"127.0.0.1"},
		"ping_interval":  []string{"100"},
		"check_interval": []string{"1"},
	}

	r.NoError(pingRunner(c, o, c.Channel()))
	get := <-c.Channel()
	r.NotNil(get)
	m, _ := get.(bogo.PingMessage)
	r.IsType(bogo.PingMessage{}, m)

	c.Cancel()
}

/* chnaged the shutdown handling gracefully
func TestPingRunner_Panic(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{
		"targets":        []string{"127.0.0.1"},
		"ping_interval":  []string{"100"},
		"check_interval": []string{"2"},
	}

	r.NoError(pingRunner(c, o, c.Channel()))
	close(c.Channel())

	c.Cancel()
}
*/

func TestPingRunner_NoTarget(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	defer c.Cancel()

	o := common.PluginOptions{
		"check_interval": []string{"1"},
	}

	err := pingRunner(c, o, c.Channel())
	r.ErrorContains(err, "no targets specified")
}

func TestPingRunner_EmptyTarget(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	defer c.Cancel()

	o := common.PluginOptions{
		"targets":        []string{""},
		"check_interval": []string{"1"},
	}

	err := pingRunner(c, o, c.Channel())
	r.ErrorContains(err, "target string should not be empty")
}

func TestPingRunner_WrongTarget(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)

	o := common.PluginOptions{
		"targets":        []string{"badaddress"},
		"ping_interval":  []string{"100"},
		"check_interval": []string{"1"},
	}

	r.NoError(pingRunner(c, o, c.Channel()))
	time.Sleep(1100 * time.Millisecond)

	c.Cancel()
}

func TestPingRunner_InvalidOptionValueCheckInterval(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	defer c.Cancel()

	o := common.PluginOptions{
		"targets":        []string{"127.0.0.1"},
		"check_interval": []string{"number"},
	}

	err := pingRunner(c, o, c.Channel())
	r.ErrorContains(err, "invalid option value")
}

func TestPingRunner_InvalidOptionValuePingInterval(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	defer c.Cancel()

	o := common.PluginOptions{
		"targets":       []string{"127.0.0.1"},
		"ping_interval": []string{"number"},
	}

	err := pingRunner(c, o, c.Channel())
	r.ErrorContains(err, "invalid option value")
}
