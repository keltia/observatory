package observatory

import (
	"net/http"
	"time"
)

// Client is used to store proxyauth & other internal state
type Client struct {
	baseurl   string
	proxyauth string
	level     int
	client    *http.Client
	timeout   time.Duration
	cache     time.Duration

	// Local cache for 5mn of last query
	last Analyze
}

// Config is for giving options to NewClient
type Config struct {
	BaseURL string
	Timeout int
	Log     int
	Cache   int
}

// Analyse is one run
type Analyze struct {
	AlgorithmVersion int `json:"algorithm_version"`

	Grade  string `json:"grade"`
	Score  int    `json:"score"`
	ScanID int    `json:"scan_id"`

	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

	State               string `json:"state"`
	StatusCode          int    `json:"status_code"`
	Hidden              bool   `json:"hidden"`
	LikelyhoodIndicator string `json:"likelyhood_indicator"`

	TestsFailed   int `json:"tests_failed"`
	TestsPassed   int `json:"tests_passed"`
	TestsQuantity int `json:"tests_quantity"`

	ResponseHeaders map[string]string `json:"response_headers"`
}

// Scan for each individual tests
type Scan struct {
	Expectation      string `json:"expectation"`
	Name             string `json:"name"`
	Output           []byte `json:"output"`
	Pass             bool   `json:"pass"`
	Result           string `json:"result"`
	ScoreDescription string `json:"score_description"`
	ScoreModifier    int    `json:"score_modifier"`
}

// ScanReport is all results from tests
type ScanReport struct {
	Tests map[string]Scan
}
