package main

import (
	"fmt"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jace-ys/mcrawl/pkg/crawler"
	"github.com/jace-ys/mcrawl/pkg/excluders"
	"github.com/jace-ys/mcrawl/pkg/fetchers"
)

var (
	startURL  = kingpin.Arg("url", "URL to start crawling from. Will only follow URLs belonging to the given URL's subdomain.").Required().URL()
	workers   = kingpin.Flag("workers", "Number of concurrent workers to use for crawling.").Default("10").Int()
	robotstxt = kingpin.Flag("robotstxt", "Respect the site's robots.txt file, if any, while crawling.").Default("true").Bool()
	debug     = kingpin.Flag("debug", "Run the web crawler in debug mode.").Default("false").Bool()
)

func main() {
	kingpin.Parse()

	fetcher := fetchers.NewLinksFetcher()
	excluder, err := excluders.NewRobotsTxtExcluder("MCrawl", *robotstxt, *startURL)
	if err != nil {
		kingpin.FatalIfError(err, "failed to initialise robots.txt excluder")
	}

	crawler := crawler.NewCrawler(fetcher, excluder, *startURL, *workers, *debug)

	start := time.Now()
	results := crawler.Crawl()
	duration := time.Now().Sub(start).Seconds()

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
	fmt.Printf("Time taken: %.3fs\n", duration)
}
