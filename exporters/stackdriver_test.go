package exporters

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/hyeoncheon/bogo"
	"github.com/hyeoncheon/bogo/internal/common"

	"github.com/stretchr/testify/require"
)

func TestRegisterStackdriver(t *testing.T) {
	r := require.New(t)

	p := (&Exporter{}).RegisterStackdriver()
	r.IsType(&Exporter{}, p)
	r.Implements((*common.Plugin)(nil), p)

	r.Equal(stackdriverExporter, p.Name())
}

func TestStackdriverRunner_NotOnGCE(t *testing.T) {
	r := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	opts := common.DefaultOptions()
	c := &DummyContext{
		Context: ctx,
		Options: opts,
		cancel:  cancel,
		ch:      make(chan interface{}),
		wg:      &sync.WaitGroup{},
		logger:  common.NewDefaultLogger("info"),
		meta: &DummyMeta{
			VarWhereAmI: common.UNKNOWN,
		},
	}
	o := common.PluginOptions{}

	r.ErrorIs(stackdriverRunner(c, o, c.Channel()), common.ErrNotOnGCE)
	c.Cancel()
}

func TestStackdriverRunner_OnGCE(t *testing.T) {
	r := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	opts := common.DefaultOptions()
	c := &DummyContext{
		Context: ctx,
		Options: opts,
		cancel:  cancel,
		ch:      make(chan interface{}),
		wg:      &sync.WaitGroup{},
		logger:  common.NewDefaultLogger("info"),
		meta: &DummyMeta{
			VarWhereAmI: common.GOOGLE,
		},
	}
	o := common.PluginOptions{}

	r.Error(stackdriverRunner(c, o, c.Channel()))
	c.Cancel()
}

func TestRecordPingMessage(t *testing.T) {
	r := require.New(t)
	r.NoError(recordPingMessage(&reporter{
		instanceName: "instance",
		externalIP:   "ipaddress",
		zone:         "here",
	}, &bogo.PingMessage{
		Addr:   "localhost",
		IPAddr: &net.IPAddr{},
		Count:  10,
		Loss:   10,
		MinRtt: 10,
		MaxRtt: 10,
		AvgRtt: 10,
		StdDev: 10,
	}))
}
