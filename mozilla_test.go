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

func TestClient_GetGrade(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr, err := ioutil.ReadFile("testdata/ssllabs-post.json")
	assert.NoError(t, err)

	ftc, err := ioutil.ReadFile("testdata/ssllabs-get.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Post("analyze").
		MatchParam("host", site).
		MatchHeaders(map[string]string{
			"content-type": "application/json",
			"accept":       "application/json",
		}).
		BodyString("hidden=true&rescan=true").
		Reply(200).
		BodyString(string(ftr))

	gock.New(baseURL).
		Get("analyze").
		MatchParam("host", site).
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	ret, err := c.GetGrade(site)
	assert.NoError(t, err)
	assert.EqualValues(t, "A+", ret)
}

func TestClient_GetScore(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr, err := ioutil.ReadFile("testdata/ssllabs-post.json")
	assert.NoError(t, err)

	ftc, err := ioutil.ReadFile("testdata/ssllabs-get.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Post("analyze").
		MatchParam("host", site).
		MatchHeaders(map[string]string{
			"content-type": "application/json",
			"accept":       "application/json",
		}).
		BodyString("hidden=true&rescan=true").
		Reply(200).
		BodyString(string(ftr))

	gock.New(baseURL).
		Get("analyze").
		MatchParam("host", site).
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	ret, err := c.GetScore(site)
	assert.NoError(t, err)
	assert.EqualValues(t, 105, ret)
}

func TestClient_GetScanID(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftc, err := ioutil.ReadFile("testdata/ssllabs-get.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("analyze").
		MatchParam("host", site).
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	ret, err := c.GetScanID(site)
	assert.NoError(t, err)
	assert.EqualValues(t, 8507653, ret)

}

func TestClient_GetScanReport(t *testing.T) {
	defer gock.Off()

	ftc, err := ioutil.ReadFile("testdata/ssllabs-8507653.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("getScanResults").
		MatchParam("scan", "8507653").
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	ret, err := c.GetScanResults(8507653)
	assert.NoError(t, err)
	assert.EqualValues(t, ftc, ret)

}
