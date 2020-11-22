package crawler

import (
	urlpkg "net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jace-ys/mcrawl/pkg/crawler/fakes"
)

var (
	startURL, _ = urlpkg.Parse("https://golang.org/")
)

func TestCrawl(t *testing.T) {
	fetcher := fakes.NewFakeFetcher()
	excluder := fakes.NewFakeExcluder(false)
	crawler := NewCrawler(fetcher, excluder, startURL, 2, false)

	results := crawler.Crawl()
	assert.Equal(t, len(results), 5)

	tt := []struct {
		target string
		links  []string
	}{
		{
			target: "https://golang.org",
			links: []string{
				"https://golang.org/pkg",
				"https://golang.org/cmd",
				"https://monzo.com",
			},
		},
		{
			target: "https://golang.org/pkg",
			links: []string{
				"https://golang.org",
				"https://golang.org/cmd",
				"https://golang.org/pkg/fmt",
				"https://golang.org/pkg/os",
			},
		},
		{
			target: "https://golang.org/pkg/fmt",
			links: []string{
				"https://golang.org",
				"https://golang.org/pkg",
			},
		},
		{
			target: "https://golang.org/pkg/os",
			links: []string{
				"https://golang.org",
				"https://golang.org/pkg",
			},
		},
		{
			target: "https://golang.org/cmd",
			links:  []string{},
		},
	}

	for _, tc := range tt {
		assert.Contains(t, results, tc.target)
		assert.NoError(t, results[tc.target].Err)
		assert.Equal(t, results[tc.target].Links, tc.links)
	}
}

func TestCrawlWithExclude(t *testing.T) {
	fetcher := fakes.NewFakeFetcher()
	excluder := fakes.NewFakeExcluder(true)
	crawler := NewCrawler(fetcher, excluder, startURL, 2, false)

	results := crawler.Crawl()
	assert.Equal(t, len(results), 0)
}

func TestShouldFollow(t *testing.T) {
	fetcher := fakes.NewFakeFetcher()
	excluder := fakes.NewFakeExcluder(false)
	crawler := NewCrawler(fetcher, excluder, startURL, 2, false)

	tt := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Same Subdomain & Scheme",
			url:      "https://golang.org/help",
			expected: true,
		},
		{
			name:     "Different Subdomain",
			url:      "https://monzo.com/",
			expected: false,
		},
		{
			name:     "Different Scheme",
			url:      "mailto:tech-hiring@monzo.com",
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := crawler.shouldFollow(tc.url)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
