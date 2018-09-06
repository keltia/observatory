package observatory

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/h2non/gock"

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

func TestClient_GetHostHistory(t *testing.T) {
	c, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)

	_, err = c.GetHostHistory("")
	assert.Error(t, err)
}

func TestClient_GetHostHistory2(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr, err := ioutil.ReadFile("testdata/ssllabs-history.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("getHostHistory").
		MatchParam("host", site).
		Reply(200).
		BodyString(string(ftr))

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var hh []HostHistory

	err = json.Unmarshal(ftr, &hh)
	require.NoError(t, err)

	h, err := c.GetHostHistory(site)
	assert.NoError(t, err)
	assert.IsType(t, ([]HostHistory)(nil), h)

	assert.EqualValues(t, hh, h)
}

func TestVersion(t *testing.T) {
	v := Version()
	require.Equal(t, MyVersion, v)
}
