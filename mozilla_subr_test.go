package observatory

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testURL = "http://127.0.0.1:10000"

func TestToDuration(t *testing.T) {
	testData := []struct {
		In  int
		Out time.Duration
	}{
		{0, 0},
		{1, 1 * time.Second},
		{50, 50 * time.Second},
	}

	for _, td := range testData {
		assert.Equal(t, td.Out, toDuration(td.In))
	}
}

func TestMyRedirect(t *testing.T) {
	err := myRedirect(nil, nil)

	assert.NoError(t, err)
}

func TestAddQueryParameters(t *testing.T) {
	p := AddQueryParameters("", map[string]string{})
	assert.Equal(t, "", p)
}

func TestAddQueryParameters_1(t *testing.T) {
	p := AddQueryParameters("", map[string]string{"": ""})
	assert.Equal(t, "?=", p)
}

func TestAddQueryParameters_2(t *testing.T) {
	p := AddQueryParameters("foo", map[string]string{"bar": "baz"})
	assert.Equal(t, "foo?bar=baz", p)
}

func Before(t *testing.T, url string) *Client {
	var testConfig = Config{BaseURL: url}

	c, err := NewClient(testConfig)
	assert.NoError(t, err)
	return c
}

func TestPrepareRequest(t *testing.T) {
	c := Before(t, testURL)
	u, _ := url.Parse(testURL)

	assert.Equal(t, testURL, c.baseurl)

	var opts = map[string]string{}
	req := c.prepareRequest("GET", "foo", opts)

	assert.IsType(t, (*http.Request)(nil), req)
	assert.Equal(t, u.Host, req.Host)
	assert.Equal(t, "GET", req.Method)
}

func TestPrepareRequest_2(t *testing.T) {
	c := Before(t, "")
	u, _ := url.Parse(baseURL)

	assert.Equal(t, baseURL, c.baseurl)

	var opts = map[string]string{}
	req := c.prepareRequest("GET", "foo", opts)

	assert.IsType(t, (*http.Request)(nil), req)
	assert.Equal(t, u.Host, req.Host)
	assert.Equal(t, "GET", req.Method)
}

func TestPrepareRequest_3(t *testing.T) {
	c := Before(t, testURL)
	u, _ := url.Parse(testURL)

	assert.Equal(t, testURL, c.baseurl)

	var opts = map[string]string{}
	req := c.prepareRequest("GET", "foo", opts)

	assert.IsType(t, (*http.Request)(nil), req)
	assert.Equal(t, u.Host, req.Host)
	assert.Equal(t, "GET", req.Method)
}

func TestPrepareRequest_4(t *testing.T) {
	c := Before(t, testURL)
	u, _ := url.Parse(testURL)

	var opts = map[string]string{}
	req := c.prepareRequest("POST", "foo", opts)

	assert.IsType(t, (*http.Request)(nil), req)
	assert.Equal(t, u.Host, req.Host)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, "application/json", req.Header.Get("Accept"))
}

func TestClient_CallAPI(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr := `{"error":"recent-scan-not-found","text":"Recently completed scan for www.ssllabs.com not found"}`

	gock.New(baseURL).
		Post("analyze").
		MatchParam("host", site).
		MatchHeaders(map[string]string{
			"content-type": "application/json",
			"accept":       "application/json",
		}).
		BodyString("hidden=true").
		Reply(200).
		BodyString(ftr)

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	opts := map[string]string{
		"host": site,
	}

	body := "hidden=true"
	ret, err := c.callAPI("POST", "analyze", body, opts)

	assert.NoError(t, err)
	assert.Equal(t, ftr, string(ret))
}

func TestClient_CallAPI2(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr, err := ioutil.ReadFile("testdata/ssllabs-post.json")
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

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	opts := map[string]string{
		"host": site,
	}

	body := "hidden=true&rescan=true"
	ret, err := c.callAPI("POST", "analyze", body, opts)

	assert.NoError(t, err)
	assert.Equal(t, ftr, ret)
}

func TestClient_CallAPI3(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr, err := ioutil.ReadFile("testdata/ssllabs-get.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("analyze").
		MatchParam("host", site).
		Reply(200).
		BodyString(string(ftr))

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	opts := map[string]string{
		"host": site,
	}

	ret, err := c.callAPI("GET", "analyze", "", opts)

	assert.NoError(t, err)
	assert.Equal(t, ftr, ret)
}

func TestClient_GetAnalyse(t *testing.T) {
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

	c, err := NewClient(Config{Timeout: 10, Log: 2})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var report Analyze

	err = json.Unmarshal(ftc, &report)
	require.NoError(t, err)

	ret, err := c.getAnalyze(site, true)
	assert.NoError(t, err)
	assert.EqualValues(t, &report, ret)
}

// It will loop & retry
func TestClient_GetAnalyse2(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftc, err := ioutil.ReadFile("testdata/ssllabs-post.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("analyze").
		MatchParam("host", site).
		Reply(200).
		BodyString(string(ftc))

	c, err := NewClient(Config{Timeout: 10, Log: 2})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var report Analyze

	err = json.Unmarshal(ftc, &report)
	require.NoError(t, err)

	raw, err := c.getAnalyze(site, false)
	assert.Error(t, err)
	t.Logf("error=%v raw=%v", err, raw)
}

func TestClient_GetAnalyse_Error(t *testing.T) {
	defer gock.Off()

	site := "www.ssllabs.com"

	ftr, err := ioutil.ReadFile("testdata/ssllabs-post.json")
	assert.NoError(t, err)

	ftc, err := ioutil.ReadFile("testdata/ssllabs-error.json")
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

	c, err := NewClient(Config{Timeout: 10, Log: 2})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var report Analyze

	err = json.Unmarshal(ftc, &report)
	require.NoError(t, err)

	ret, err := c.getAnalyze(site, true)
	assert.Error(t, err)
	assert.EqualValues(t, &report, ret)
}

func TestClient_GetAnalyseEmpty(t *testing.T) {
	defer gock.Off()

	site := ""

	c, err := NewClient(Config{Timeout: 10})
	assert.NoError(t, err)
	assert.Equal(t, baseURL, c.baseurl)

	_, err = c.getAnalyze(site, false)
	assert.Error(t, err)
	assert.Equal(t, "empty site", err.Error())
}

func TestClient_GetScanReport2(t *testing.T) {

}

func TestIsValid_Nil(t *testing.T) {
	require.False(t, isValid(nil))
}

func TestIsValid_Old(t *testing.T) {
	now := time.Now().Add(-1 * time.Minute)
	ar := &Analyze{EndTime: now.Format(time.RFC1123)}
	require.True(t, isValid(ar))
}

func TestIsValid_New(t *testing.T) {
	now := time.Now().Add(-20 * time.Minute)
	ar := &Analyze{EndTime: now.Format(time.RFC1123)}
	require.False(t, isValid(ar))
}
