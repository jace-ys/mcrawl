package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jace-ys/mcrawl/pkg/crawler"
)

var (
	startURL = kingpin.Arg("url", "URL to start crawling from. Will only follow URLs belonging to the given URL's subdomain.").Required().URL()
	workers  = kingpin.Flag("workers", "Number of concurrent workers to use for crawling.").Default("10").Int()
	debug    = kingpin.Flag("debug", "Run the web crawler in debug mode.").Default("false").Bool()
)

func main() {
	kingpin.Parse()

	crawler := crawler.NewCrawler(crawler.NewFakeFetcher(), *startURL, *workers, *debug)

	results := crawler.Crawl()

	success := 0
	for target, result := range results {
		if result.Err == nil {
			fmt.Printf("%s\n", target)
			for _, link := range result.Links {
				fmt.Printf("  -> %s\n", link)
			}

			success++
		}
	}

	fmt.Println("======================")
	fmt.Printf("Unique URLs crawled: %d\n", success)
}
