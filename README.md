observatory
==========

[![GitHub release](https://img.shields.io/github/release/keltia/observatory.svg)](https://github.com/keltia/observatory/releases)
[![GitHub issues](https://img.shields.io/github/issues/keltia/observatory.svg)](https://github.com/keltia/observatory/issues)
[![Go Version](https://img.shields.io/badge/go-1.10-blue.svg)](https://golang.org/dl/)
[![Build Status](https://travis-ci.org/keltia/observatory.svg?branch=master)](https://travis-ci.org/keltia/observatory)
[![GoDoc](http://godoc.org/github.com/keltia/observatory?status.svg)](http://godoc.org/github.com/keltia/observatory)
[![SemVer](http://img.shields.io/SemVer/2.0.0.png)](https://semver.org/spec/v2.0.0.html)
[![License](https://img.shields.io/pypi/l/Django.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/observatory)](https://goreportcard.com/report/github.com/keltia/observatory)

Go wrapper for [Mozilla Observatory](https://observatory.mozilla.org/) API.

## Requirements

* Go >= 1.10

`github.com/keltia/observatory` is a Go module (you can use either Go 1.10 with `vgo` or 1.11+).  The API exposed follows the Semantic Versioning scheme to guarantee a consistent API compatibility.

## Installation

You need to install my `proxy` module before if you are using Go 1.10.x or earlier.

    go get github.com/keltia/proxy

With Go 1.11+ and its modules support, it should work out of the box with

    go get github.com/keltia/observatory/cmd/...

if you have the `GO111MODULE` environment variable set on `on`.

## CLI

There is a small example program included in `cmd/observatory` to either show the grade of a given site or JSON dump of the detailed report.

Easy to use:
```
    $ observatory www.ssllabs.com
    observatory Wrapper: 0.3.0 API version 1.2.0
    
    Grade for 'www.ssllabs.com' is A+
```

You can use [`jq`](https://stedolan.github.io/jq/) to display the output of `observatory -d <site>` in a colorised way:

    observatory -d observatory.mozilla.org | jq .

## API Usage

As with many API wrappers, you will need to first create a client with some optional configuration, then there are two main functions:

``` go
    // Simplest way
    c, _ := observatory.NewClient()
    grade, err := c.GetScore("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
```

If you want to change the default options, you need to create a `ssllabs.Config` object and pass it to `NewClient`:

``` go
    // With some options, timeout at 15s, caching for 10s and debug-like verbosity
    cnf := observatory.Config{
        Timeout:15,
        Retries:3,
        Log:2,
    }
    c, err := observatory.NewClient(cnf)
    report, err := c.GetScore("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
```

OPTIONS for NewClient()

| Option  | Type | Description |
| ------- | ---- | ----------- |
| Timeout | int  | time for connections (default: 10s) |
| Log     | int  | 1: verbose, 2: debug (default: 0) |
| Retries | int  | Number of retries when not FINISHED (default: 5) |
| Refresh | bool | Force refresh of the sites (default: false) |

For the `GetScanResults()` call, the raw JSON object will be returned (and presumably handled by `jq`).

``` go
    // Simplest way
    c, _ := observatory.NewClient()
    
    scanid, err := c.GetScanID("example.com")
    
    report, err := c.GetScanResults(scanid)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    fmt.Printf("Full report:\n%v\n", report)
```

The `GetHostHistory()` returns the list of recent scans for the given site:

``` go
    // Simplest way
    c, _ := observatory.NewClient()
    
    scans, err := c.GetHostHistory("example.com")
    for _, s := range scans {
        ...
    }
```

There is no top-level `GetGrade` function but it is very easy to implement:

``` go
    func GetGrade(site string) string {
        g, _ := observatory.NewClient().GetGrade(site)
        return g
    }
```

### NOTE

v1.1.x implemented the `GetScanReport` call but that does not correspond to any real API calls.  It is now just an alias to `GetScanResults`.  DO NOT USE IT.  DEPRECATED.

### API Calls Implemented

- `analyze`
- `getScanResults`
- `getHostHistory`

### API NOT Implemented

- `getRecentScans`

## Using behind a web Proxy

Dependency: proxy support is provided by my `github.com/keltia/proxy` module.

UNIX/Linux:

```
    export HTTP_PROXY=[http://]host[:port] (sh/bash/zsh)
    setenv HTTP_PROXY [http://]host[:port] (csh/tcsh)
```

Windows:

```
    set HTTP_PROXY=[http://]host[:port]
```

The rules of Go's `ProxyFromEnvironment` apply (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, lowercase variants allowed).

If your proxy requires you to authenticate, please create a file named `.netrc` in your HOME directory with permissions either `0400` or `0600` with the following data:

    machine proxy user <username> password <password>

and it should be picked up. On Windows, the file will be located at

    %LOCALAPPDATA%\observatory\netrc

## License

The [BSD 2-Clause license](https://github.com/keltia/observatory/LICENSE.md).

# Contributing

This project is an open Open Source project, please read `CONTRIBUTING.md`.

# References

[Mozilla Observatory documentation](https://github.com/mozilla/http-observatory/blob/master/httpobs/docs/api.md#host-history)

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.