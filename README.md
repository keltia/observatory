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

## USAGE

There is a small example program included in `cmd/getgrade` to either show the grade of a given site or JSON dump of the detailed report.

You can use `jq` to display the output of `getgrade -d <site>` in a colorised way:

    getgrade -d observatory.mozilla.org | jq .

## API Usage

As with many API wrappers, you will need to first create a client with some optional configuration, then there are two main functions:

``` go
    // Simplest way
    c := observatory.NewClient()
    grade, err := c.GetScore("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }


    // With some options, timeout at 15s and debug-like verbosity
    cnf := observatory.Config{
        Timeout:15,
        Log:2,
    }
    c := observatory.NewClient(cnf)
    report, err := c.GetDetailedReport("foo.xxx")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
```

OPTIONS

| Option  | Type | Description |
| ------- | ---- | ----------- |
| Timeout | int  | time for connections (default: 10s) |
| Log     | int  | 1: verbose, 2: debug (default: 0) |
| Refresh | bool | Force refresh of the sites (default: false) |
| Cache   | int  | time allowed for caching last call (default 300s) |


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

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.