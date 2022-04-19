package meari

import (
	"testing"

	"github.com/hyeoncheon/bogo/internal/common"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	r := require.New(t)

	opts := common.DefaultOptions()
	ctx, _ := common.NewDefaultContext(&opts)

	s, err := NewServer(ctx, &opts)
	r.NoError(err)
	r.NotNil(s)
}
