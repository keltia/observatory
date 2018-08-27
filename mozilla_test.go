package observatory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient()

	assert.NoError(t, err)
	assert.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
}
