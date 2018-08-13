package observatory

/*
Not going to implement the full scan report struct, I do not need it, juste grade/score
*/
import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/keltia/proxy"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://http-observatory.security.mozilla.org/api/v1"

	// DefaultWait is the timeout
	DefaultWait = 10 * time.Second

	// DefaultCache is in second
	DefaultCache = 5 * time.Minute

	// MyVersion is the API version
	MyVersion = "0.3.0"

	// MyName is the name used for the configuration
	MyName = "observatory"
)

var (
	APIVersion = MyVersion
)

// Public functions

// NewClient setups proxy authentication
func NewClient(cnf ...Config) (*Client, error) {
	var c *Client

	// Set default
	if len(cnf) == 0 {
		c = &Client{
			baseurl: baseURL,
			timeout: DefaultWait,
			cache:   DefaultCache,
		}
	} else {
		c = &Client{
			baseurl: cnf[0].BaseURL,
			level:   cnf[0].Log,
			cache:   toDuration(cnf[0].Cache) * time.Second,
			timeout: toDuration(cnf[0].Timeout) * time.Second,
		}

		if cnf[0].Timeout == 0 {
			c.timeout = DefaultWait
		} else {
			c.timeout = time.Duration(cnf[0].Timeout) * time.Second
		}

		// Ensure we have the API endpoint right
		if c.baseurl == "" {
			c.baseurl = baseURL
		}

		c.debug("got cnf: %#v", cnf[0])
	}

	// We do not care whether it fails or not, if it does, just no proxyauth.
	proxyauth, _ := proxy.SetupProxyAuth()

	// Save it
	c.proxyauth = proxyauth
	c.debug("got proxyauth: %s", c.proxyauth)

	_, trsp := proxy.SetupTransport(c.baseurl)
	c.client = &http.Client{
		Transport:     trsp,
		Timeout:       c.timeout,
		CheckRedirect: myRedirect,
	}
	c.debug("mozilla: c=%#v", c)
	return c, nil
}

// GetScore returns the integer value of the grade
func (c *Client) GetScore(site string) (score int, err error) {
	c.debug("GetScore")

	// Check whether we have a cached value inside our caching timeout

	opts := map[string]string{
		"host": site,
	}

	_, err = c.callAPI("POST", "analyze", "hidden=true&rescan=true", opts)
	if err != nil {
		return -1, errors.Wrap(err, "callAPI failed")
	}
	r, err := c.callAPI("GET", "analyze", "", opts)

	var ar Analyze

	err = json.Unmarshal(r, &ar)
	return ar.Score, errors.Wrap(err, "GetScore failed")
}

// GetGrade returns the letter equivalent to the score
func (c *Client) GetGrade(site string) (grade string, err error) {
	c.debug("GetGrade")

	opts := map[string]string{
		"host": site,
	}

	_, err = c.callAPI("POST", "analyze", "hidden=true&rescan=true", opts)
	if err != nil {
		return "Z", errors.Wrap(err, "callAPI failed")
	}
	r, err := c.callAPI("GET", "analyze", "", opts)
	if err != nil {
		return "Z", errors.Wrap(err, "callAPI failed")
	}

	var ar Analyze

	err = json.Unmarshal(r, &ar)
	return ar.Grade, errors.Wrap(err, "GetGrade failed")
}

// GetDetailedReport returns the full scan report
func (c *Client) GetScanReport(scanID string) (ScanReport, error) {
	c.debug("GetScanReport")

	opts := map[string]string{
		"scan": scanID,
	}

	s, err := c.callAPI("GET", "getScanResults", "", opts)

	var sc ScanReport

	err = json.Unmarshal(s, &sc)
	return sc, errors.Wrap(err, "ScanReport unmarshall failed")
}

// Version returns guess what?
func Version() string {
	return APIVersion
}
