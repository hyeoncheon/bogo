package exporters

import (
	"context"
	"sync"

	"github.com/hyeoncheon/bogo/internal/common"
)

type DummyContext struct {
	context.Context // nolint
	common.Options
	cancel context.CancelFunc
	ch     chan interface{}
	wg     *sync.WaitGroup
	logger common.Logger
	meta   common.MetaClient
}

var _ common.Context = &DummyContext{}

// Cancel implements common.Context.
func (c *DummyContext) Cancel() {
	c.cancel()
	c.wg.Wait()
	close(c.ch)
}

// Logger implements common.Context.
func (c *DummyContext) Logger() common.Logger {
	return c.logger
}

// Meta implements common.Context.
func (c *DummyContext) Meta() common.MetaClient {
	return c.meta
}

// WG implements common.Context.
func (c *DummyContext) WG() *sync.WaitGroup {
	return c.wg
}

// Channel implements common.Context.
func (c *DummyContext) Channel() chan interface{} {
	return c.ch
}

type DummyMeta struct {
	VarExternalIP   string
	VarInstanceName string
	VarWhereAmI     string
	VarZone         string
}

// AttributeCSV implements common.MetaClient.
func (*DummyMeta) AttributeCSV(_ string) []string {
	panic("unimplemented")
}

// AttributeSSV implements common.MetaClient.
func (*DummyMeta) AttributeSSV(_ string) []string {
	panic("unimplemented")
}

// AttributeValue implements common.MetaClient.
func (*DummyMeta) AttributeValue(_ string) string {
	panic("unimplemented")
}

// AttributeValues implements common.MetaClient.
func (*DummyMeta) AttributeValues(_ string) []string {
	panic("unimplemented")
}

// ExternalIP implements common.MetaClient.
func (m *DummyMeta) ExternalIP() string {
	return m.VarExternalIP
}

// InstanceName implements common.MetaClient.
func (m *DummyMeta) InstanceName() string {
	return m.VarInstanceName
}

// WhereAmI implements common.MetaClient.
func (m *DummyMeta) WhereAmI() string {
	return m.VarWhereAmI
}

// Zone implements common.MetaClient.
func (m *DummyMeta) Zone() string {
	return m.VarZone
}
