package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	startURL = kingpin.Arg("url", "URL to start crawling from. Will only follow URLs belonging to the given URL's subdomain.").Required().URL()
	workers  = kingpin.Flag("workers", "Number of concurrent workers to use for crawling.").Default("10").Int()
	debug    = kingpin.Flag("debug", "Run the web crawler in debug mode.").Default("false").Bool()
)

func main() {
	kingpin.Parse()
}
