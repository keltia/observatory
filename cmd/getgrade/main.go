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
	"github.com/keltia/observatory"
	"log"
	"os"
	"path/filepath"
)

var (
	fDetailed bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}
}

func main() {
	site := flag.Arg(0)

	// Setup client
	c, err := observatory.NewClient()
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
			MyName, observatory.MyVersion, observatory.APIVersion)
		grade, err := c.GetScore(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s'\n", site)
		}
		fmt.Printf("Grade for '%s' is %s\n", site, grade)
	}
}
