package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAllHandlers tests if AllHandlers works correctly.
func TestAllHandlers(t *testing.T) {
	r := require.New(t)

	handlers := AllHandlers()
	r.NotNil(handlers)
	// this will not be a fixed number, and not a good test point but...
	r.Equal(2, len(handlers))
}
