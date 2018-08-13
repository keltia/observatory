package observatory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// prepareRequest insert all pre-defined stuff
func (c *Client) prepareRequest(method, what string, opts map[string]string) (req *http.Request) {
	var endPoint string

	// This allow for overriding baseurl for tests
	if c.baseurl != "" {
		endPoint = fmt.Sprintf("%s/%s", c.baseurl, what)
	} else {
		endPoint = fmt.Sprintf("%s/%s", baseURL, what)
	}

	c.verbose("Options:\n%v", opts)
	baseURL := AddQueryParameters(endPoint, opts)
	c.debug("baseURL: %s", baseURL)

	req, _ = http.NewRequest(method, baseURL, nil)

	c.debug("req=%#v", req)

	// We need these when we POST
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	return
}

func (c *Client) callAPI(word, cmd, sbody string, opts map[string]string) ([]byte, error) {

	req := c.prepareRequest(word, cmd, opts)
	if req == nil {
		return []byte{}, errors.New("req is nil")
	}

	c.debug("req=%#v", req)
	c.debug("clt=%#v", c.client)
	c.debug("opts=%v", opts)

	// If we have a POST and a body, insert them.
	if sbody != "" && word == "POST" {
		body := []byte(sbody)
		buf := bytes.NewReader(body)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(buf.Len())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.debug("err=%#v", err)
		return []byte{}, errors.Wrap(err, "1st call failed")
	}
	c.debug("resp=%#v", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.Wrap(err, "can not read body")
	}

	c.debug("body=%v", string(body))

	if resp.StatusCode == http.StatusOK {

		c.debug("status OK")

		if string(body) == "pending" {
			time.Sleep(10 * time.Second)
			resp, err = c.client.Do(req)
			if err != nil {
				return body, errors.Wrap(err, "pending failed")
			}
			c.debug("resp was %v", resp)
		}
	} else if resp.StatusCode == http.StatusFound {
		str := resp.Header["Location"][0]

		c.debug("Got 302 to %s", str)

		req := c.prepareRequest(word, cmd, opts)
		if err != nil {
			return []byte{}, errors.Wrap(err, "Cannot handle redirect")
		}

		resp, err = c.client.Do(req)
		if err != nil {
			return []byte{}, errors.Wrap(err, "client.Do failed")
		}
		c.debug("resp was %v", resp)
	} else {
		return body, errors.Wrapf(err, "bad status code: %v body: %q", resp.Status, body)
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
			return nil, errors.Wrapf(err, "getAnalyze - ret: %v", ret)
		}
	}

	r, err := c.callAPI("GET", "analyze", "", opts)
	if err != nil {
		return &ar, errors.Wrap(err, "getAnalyze")
	}

	err = json.Unmarshal(r, &ar)
	return &ar, errors.Wrap(err, "getAnalyze")
}
