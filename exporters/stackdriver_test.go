package exporters

import (
	"testing"

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
