package observatory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Private area

func toDuration(t int) time.Duration {
	s := fmt.Sprintf("%ds", t)
	d, _ := time.ParseDuration(s)
	return d
}

func myRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

// AddQueryParameters adds query parameters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	params := url.Values{}
	if len(queryParams) == 0 {
		return baseURL
	}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// prepareRequest insert all pre-defined stuff
func (c *Client) prepareRequest(method, what string, opts map[string]string) (req *http.Request) {
	var endPoint string

	endPoint = fmt.Sprintf("%s/%s", c.baseurl, what)

	c.verbose("Options:\n%v", opts)
	baseURL := AddQueryParameters(endPoint, opts)
	c.debug("baseURL: %s", baseURL)

	req, _ = http.NewRequest(method, baseURL, nil)

	// We need these when we POST
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	return
}

// callAPI is the main API call — straightforward, clean logic
func (c *Client) callAPI(word, cmd, sbody string, opts map[string]string) ([]byte, error) {
	c.debug("callAPI")
	req := c.prepareRequest(word, cmd, opts)
	if req == nil {
		return []byte{}, errors.New("req is nil")
	}

	c.debug("clt=%#v", c.client)
	c.debug("opts=%v", opts)

	// If we have a POST and a body, insert them.
	if sbody != "" && word == "POST" {
		buf := bytes.NewBufferString(sbody)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(buf.Len())
	}

	c.debug("req=%#v body=%v", req, req.Body)

	resp, err := c.client.Do(req)
	if err != nil {
		c.debug("err=%#v", err)
		return []byte{}, errors.Wrap(err, "1st call")
	}
	defer resp.Body.Close()

	c.debug("resp=%#v", resp)

	c.debug("read body")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.Wrap(err, "body read")
	}

	c.debug("body=%v", string(body))

	if resp.StatusCode == http.StatusOK {

		c.debug("status OK")

		if strings.Contains(string(body), "error:") {
			return body, errors.New("error")
		}
	} else {
		return body, errors.Wrapf(err, "status: %v body: %q", resp.Status, body)
	}
	return body, err
}

// getAnalyze is an helper func for the API
func (c *Client) getAnalyze(site string, force bool) (*Analyze, error) {
	var ar Analyze

	opts := map[string]string{
		"host": site,
	}

	body := "hidden=true"
	if force {
		body = body + "&rescan=true"
		ret, err := c.callAPI("POST", "analyze", body, opts)
		if err != nil {
			return &Analyze{}, errors.Wrapf(err, "getAnalyze - POST: %s", string(ret))
		}
	}

	r, err := c.callAPI("GET", "analyze", "", opts)
	if err != nil {
		return &ar, errors.Wrap(err, "getAnalyze - GET")
	}

	err = json.Unmarshal(r, &ar)
	return &ar, errors.Wrapf(err, "getAnalyze: %#v", r)
}
