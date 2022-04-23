package meari

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

// TestNewServer tests NewServer with default Options and default Context.
func TestNewServer(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	ctx, _ := common.NewDefaultContext(&opts)

	s, err := NewServer(ctx, &opts)
	r.NoError(err)
	r.NotNil(s)
}
