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

func BeforeAPI(t *testing.T) {
	/*	var err error


		// define request->response pairs
		request1, _ := url.Parse(testURL + "/analyze?host=lbl.gov")
		request2, _ := url.Parse(testURL + "/getScanResults?scan=8442544")

		t.Logf("r1=%s", request1.String())
		ftq, err = ioutil.ReadFile("testdata/lbl.gov.json")
		assert.NoError(t, err)

		ftr, err = ioutil.ReadFile("testdata/lbl.gov.data.json")
		assert.NoError(t, err)

		aresp := []httpmock.MockResponse{
			{
				Request: http.Request{
					Method: "POST",
					URL:    request1,
					Header: map[string][]string{
						"content-type": {"application/json"},
						"accept":       {"application/json"},
					},
					ContentLength: int64(len("hidden=true&rescan=true")),
					Body:          ioutil.NopCloser(strings.NewReader("hidden=true&rescan=true")),
				},
				Response: httpmock.Response{
					StatusCode: 200,
					Body:       string(ftq),
				},
			},
			{
				Request: http.Request{
					Method: "GET",
					URL:    request2,
				},
				Response: httpmock.Response{
					StatusCode: 200,
					Body:       string(ftr),
				},
			},
		}

		mockService.AddResponses(aresp)
		t.Logf("respmap=%v", mockService.ResponseMap)
	*/
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
