package observatory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	cd = Config{
		BaseURL: baseURL,
		Retries: DefaultRetry,
	}
)

func TestNewClient(t *testing.T) {
	c, err := NewClient()

	assert.NoError(t, err)
	assert.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)

	assert.Equal(t, cd.BaseURL, c.baseurl)
	assert.Equal(t, cd.Retries, c.retries)
	assert.Equal(t, 0, c.level)

	assert.NotNil(t, c.client)
}
