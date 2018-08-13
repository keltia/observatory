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

	/*str := fmt.Sprintf("%s/%s?host=%s", c.baseurl, cmd, site)

	c.debug("str=%s", str)
	req, err := http.NewRequest(word, str, nil)
	*/
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
		c.verbose("err=%#v", err)
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
			c.verbose("resp was %v", resp)
		}
	} else if resp.StatusCode == http.StatusFound {
		str := resp.Header["Location"][0]

		c.debug("Got 302 to %s", str)

		req, err = http.NewRequest("GET", str, nil)
		if err != nil {
			return []byte{}, errors.Wrap(err, "Cannot handle redirect")
		}

		resp, err = c.client.Do(req)
		if err != nil {
			return []byte{}, errors.Wrap(err, "client.Do failed")
		}
		c.verbose("resp was %v", resp)
	} else {
		return body, errors.Wrapf(err, "bad status code: %v body: %q", resp.Status, body)
	}

	var report Analyze

	err = json.Unmarshal(body, &report)

	// Give some time to performe the test
	if report.State == "PENDING" {
		time.Sleep(2 * time.Second)
		err = nil
	}

	return body, err
}

// Mon Jan 2 15:04:05 MST 2006

func (c *Client) newEnough(endTime string) bool {
	t1, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", endTime)
	if err == nil {
		if time.Since(t1) < c.cache {
			return true
		}
	}
	return false
}
