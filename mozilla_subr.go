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

func (c *Client) callAPI(word, cmd, sbody string, opts map[string]string) ([]byte, error) {
	var body []byte

	retry := 0

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
		return body, errors.Wrap(err, "1st call")
	}
	defer resp.Body.Close()

	c.debug("resp=%#v", resp)

	for {
		if retry == c.retries {
			return nil, errors.New("retries")
		}

		c.debug("read body")
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, errors.Wrapf(err, "body read, retry=%d", retry)
		}

		c.debug("body=%v", string(body))

		if resp.StatusCode == http.StatusOK {

			c.debug("status OK")

			if strings.Contains(string(body), "error:") {
				return body, errors.New("error")
			}

			// We wait for FINISHED state
			if !strings.Contains(string(body), "FINISHED") {
				time.Sleep(2 * time.Second)
				retry++
				resp, err = c.client.Do(req)
				if err != nil {
					return body, errors.Wrapf(err, "pending, retry=%d", retry)
				}
				c.debug("resp was %v", resp)
			} else {
				return body, nil
			}
		} else {
			return body, errors.Wrapf(err, "status: %v body: %q", resp.Status, body)
		}
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
			return &Analyze{}, errors.Wrapf(err, "getAnalyze - POST: %v", ret)
		}
	}

	r, err := c.callAPI("GET", "analyze", "", opts)
	if err != nil {
		return &ar, errors.Wrap(err, "getAnalyze - GET")
	}

	err = json.Unmarshal(r, &ar)
	return &ar, errors.Wrapf(err, "getAnalyze: %#v", r)
}
