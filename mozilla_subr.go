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

	c.debug("Options:\n%v", opts)
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

func isValid(ar *Analyze) bool {
	if ar == nil {
		return false
	}

	now := time.Now()
	last, _ := time.Parse(time.RFC1123, ar.EndTime)
	last = last.Add(10 * time.Minute)
	return last.After(now)
}

// getAnalyze is an helper func for the API — where the loop/waiting appears
func (c *Client) getAnalyze(site string, force bool) (*Analyze, error) {
	var (
		raw []byte
		ar  Analyze
	)

	if site == "" {
		return &Analyze{}, fmt.Errorf("empty site")
	}

	opts := map[string]string{
		"host": site,
	}

	if force {
		body := "hidden=true&rescan=true"
		ret, err := c.callAPI("POST", "analyze", body, opts)
		if err != nil || strings.Contains(string(ret), `"error":`) {
			c.debug("post/1st call")
			return &Analyze{}, errors.Wrapf(err, "post/Analyze: %s", string(ret))
		}
	}

	retry := 0

	// Cached value is usable?
	if isValid(c.last) {
		return c.last, nil
	}

	// WAIT/RETRY loop is only for Analyse.
	for {
		if retry >= c.retries {
			c.debug("too many retries")
			return &Analyze{}, fmt.Errorf("retries exceeded - raw=%v", raw)
		}

		raw, err := c.callAPI("GET", "analyze", "", opts)
		if err != nil {
			c.debug("get/analyse")
			return &ar, errors.Wrapf(err, "get/Analyze: %v", raw)
		}

		if strings.Contains(string(raw), `state":"PENDING"`) {
			c.debug("PENDING retry=%d", retry)
			time.Sleep(2 * time.Second)
			retry++
			continue
		}

		if strings.Contains(string(raw), `"state":"FAILED"`) {
			c.debug("FAILED retry=%d", retry)
			c.debug("raw/analyse=%s", string(raw))

			_ = json.Unmarshal(raw, &ar)
			return &ar, errors.New("site analysis failed")
		}

		if strings.Contains(string(raw), `state":"FINISHED"`) {
			c.debug("FINISHED retry=%d", retry)
			c.debug("raw/analyse=%s", string(raw))

			err := json.Unmarshal(raw, &ar)
			// Store the last call
			c.last = &ar
			return &ar, errors.Wrap(err, "unmarshall")
		}
		c.debug("loop retry=%d", retry)
	}
}
