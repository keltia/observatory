package observatory

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/goware/httpmock"
	"github.com/stretchr/testify/assert"
)

const testURL = "http://localhost:10000"

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
