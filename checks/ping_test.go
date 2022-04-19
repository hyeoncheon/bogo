package checks

import (
	"testing"

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
	c, cancel := common.NewDefaultContext(&opts)
	ch := make(chan interface{})
	defer close(ch)

	o := common.PluginOptions{
		"targets":        []string{"127.0.0.1"},
		"ping_interval":  []string{"100"},
		"check_interval": []string{"1"},
	}

	r.NoError(pingRunner(c, o, ch))
	get := <-ch
	r.NotNil(get)
	m := get.(bogo.PingMessage)
	r.IsType(bogo.PingMessage{}, m)

	cancel()
	c.WG().Wait()
}

func TestPingRunner_Panic(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, cancel := common.NewDefaultContext(&opts)
	defer cancel()
	ch := make(chan interface{})

	o := common.PluginOptions{
		"targets":        []string{"127.0.0.1"},
		"ping_interval":  []string{"100"},
		"check_interval": []string{"2"},
	}

	r.NoError(pingRunner(c, o, ch))
	close(ch)

	c.WG().Wait()
}

func TestPingRunner_NoTarget(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, cancel := common.NewDefaultContext(&opts)
	defer cancel()
	ch := make(chan interface{})
	defer close(ch)

	o := common.PluginOptions{
		"check_interval": []string{"1"},
	}

	err := pingRunner(c, o, ch)
	r.ErrorContains(err, "no targets specified")
}

func TestPingRunner_EmptyTarget(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, cancel := common.NewDefaultContext(&opts)
	defer cancel()
	ch := make(chan interface{})
	defer close(ch)

	o := common.PluginOptions{
		"targets":        []string{""},
		"check_interval": []string{"1"},
	}

	err := pingRunner(c, o, ch)
	r.ErrorContains(err, "target string should not be empty")
}

func TestPingRunner_InvalidOptionValueCheckInterval(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	ch := make(chan interface{})

	o := common.PluginOptions{
		"targets":        []string{"127.0.0.1"},
		"check_interval": []string{"number"},
	}

	err := pingRunner(c, o, ch)
	r.ErrorContains(err, "invalid option value")
}

func TestPingRunner_InvalidOptionValuePingInterval(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	c, _ := common.NewDefaultContext(&opts)
	ch := make(chan interface{})

	o := common.PluginOptions{
		"targets":       []string{"127.0.0.1"},
		"ping_interval": []string{"number"},
	}

	err := pingRunner(c, o, ch)
	r.ErrorContains(err, "invalid option value")
}
