package crawler

import (
	"fmt"
	urlpkg "net/url"
	"strings"
	"sync"
)

type Fetcher interface {
	// Fetch takes in a URL to be crawled and should return all unique resolved links found
	Fetch(url string) (links []string, err error)
}

type Excluder interface {
	// Exclude verifies if a path should be excluded from crawling
	Exclude(path string) bool
}

type Crawler struct {
	fetcher  Fetcher
	excluder Excluder

	startURL *urlpkg.URL
	workers  int
	debug    bool

	queueChan   chan string
	resultsChan chan Result
	wg          sync.WaitGroup
}

func NewCrawler(fetcher Fetcher, excluder Excluder, startURL *urlpkg.URL, workers int, debug bool) *Crawler {
	// Strip any trailing slash from our starting URL
	startURL.Path = strings.TrimSuffix(startURL.Path, "/")

	return &Crawler{
		fetcher:     fetcher,
		excluder:    excluder,
		startURL:    startURL,
		workers:     workers,
		debug:       debug,
		queueChan:   make(chan string, workers),
		resultsChan: make(chan Result, workers),
	}
}

// Crawl starts the crawler and returns a map containing the links found for each URL crawled
func (c *Crawler) Crawl() map[string]Result {
	results := makeResultsMap()

	// Start a number of concurrent workers for crawling
	for i := 0; i < c.workers; i++ {
		go c.crawl(i)
	}

	// Process results from our workers and push them into the results map
	go c.process(results)

	// Enqueue the URL to start crawling from
	c.enqueue(c.startURL.String())

	// Block until there are no more URLs to crawl
	c.wg.Wait()
	close(c.queueChan)
	close(c.resultsChan)

	return results.urls
}

// crawl is a worker unit used for concurrent crawling
func (c *Crawler) crawl(i int) {
	// Pick up target URLs to crawl from the queue channel
	for target := range c.queueChan {
		if c.debug {
			fmt.Printf("[worker %d] crawling: %s\n", i, target)
		}

		// Fetch links found from crawling the target URL
		links, err := c.fetcher.Fetch(target)
		if err != nil {
			// Drop any URLs that we are unable to successfully crawl
			if c.debug {
				fmt.Printf("[worker %d] error: %s\n", i, err)
			}
		}

		// Add the crawled URL and links found to the results channel
		c.resultsChan <- Result{target, links, err}
	}
}

// process reads from the results channel and spawns goroutines to enqueue additional work
func (c *Crawler) process(results *resultsMap) {
	for result := range c.resultsChan {
		// Add the result to the results map
		results.add(result)

		// Enqueue additional work asynchronously so we don't block the results channel
		go func(result Result) {
			for _, link := range result.Links {
				// Enqueue the link if we haven't crawled it yet
				if !results.contains(link) {
					c.enqueue(link)
				}
			}

			c.wg.Done()
		}(result)
	}
}

// enqueue adds the given url to the queue channel if it should be followed
func (c *Crawler) enqueue(url string) {
	if c.shouldFollow(url) {
		c.queueChan <- url
		c.wg.Add(1)
	}
}

// shouldFollow returns a bool depending on whether the given URL should be crawled
func (c *Crawler) shouldFollow(url string) bool {
	target, err := urlpkg.Parse(url)
	if err != nil {
		return false
	}

	// Only follow URLs that are part of the same subdomain as our starting URL
	if target.Host != c.startURL.Host {
		return false
	}

	// Only follow URLs that have the same scheme as our starting URL
	if target.Scheme != c.startURL.Scheme {
		return false
	}

	// Test our path against the excluder
	return !c.excluder.Exclude(target.Path)
}
