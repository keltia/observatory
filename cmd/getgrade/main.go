// main.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

/*
This is just a very short example.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/keltia/observatory"
)

var (
	fDebug    bool
	fDetailed bool
	fVerbose  bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}
}

func main() {
	var level = 0

	site := flag.Arg(0)

	if fVerbose {
		level = 1
	}

	if fDebug {
		level = 2
		fVerbose = true
	}

	// Setup client
	c, err := observatory.NewClient(observatory.Config{Log: level})
	if err != nil {
		log.Fatalf("error setting up client: %v", err)
	}

	if fDetailed {
		report, err := c.GetScanReport(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s'\n", site)
		}

		// Just dump the json
		jr, err := json.Marshal(report)
		fmt.Printf("%s\n", jr)
	} else {
		fmt.Printf("%s Wrapper: %s API version %s\n\n",
			MyName, observatory.Version(), observatory.Version())
		grade, err := c.GetGrade(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s': %v\n", site, err)
		}
		fmt.Printf("Grade for '%s' is %s\n", site, grade)
	}
}
