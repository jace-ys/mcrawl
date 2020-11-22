package fetchers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	urlpkg "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrInvalidResponseCode = fmt.Errorf("invalid response code")
	ErrInvalidContentType  = fmt.Errorf("invalid content type")
)

type LinksFetcher struct {
}

func NewLinksFetcher() *LinksFetcher {
	return &LinksFetcher{}
}

// Fetch returns all unique resolved links found in the HTML content of the given URL
func (f *LinksFetcher) Fetch(url string) ([]string, error) {
	refURL, err := urlpkg.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target URL: %w", err)
	}

	// Fetch the HTML content of the given URL
	page, err := f.fetchHTML(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HTML content: %w", err)
	}

	// Parse HTML content to find links
	found, err := f.findLinks(page)
	if err != nil {
		return nil, fmt.Errorf("failed to find links: %w", err)
	}

	// Use a map to store unique URLs after sanitising
	unique := make(map[string]struct{})
	for _, link := range found {
		// Sanitise each link using the page's URL as reference
		sanitise := f.sanitise(link, refURL)

		// Drop any links that cannot be properly sanitised
		if sanitise != "" {
			unique[sanitise] = struct{}{}
		}
	}

	var urls []string
	for url := range unique {
		urls = append(urls, url)
	}

	return urls, nil
}

// fetchHTML fetches the HTML content of the given URL and returns it as a slice of bytes
func (f *LinksFetcher) fetchHTML(url string) ([]byte, error) {
	// Make a GET request to the given URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle unsuccessful response codes
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidResponseCode, resp.StatusCode)
	}

	// Only handle valid content types for HTML
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("%w: %s", ErrInvalidContentType, contentType)
	}

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// findLinks finds all unique resolved links in the given HTML content
func (f *LinksFetcher) findLinks(data []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var found []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			found = append(found, href)
		}
	})

	return found, nil
}

// sanitise resolves and normalises the given link using the reference URL
func (f *LinksFetcher) sanitise(link string, refURL *urlpkg.URL) string {
	url, err := urlpkg.Parse(link)
	if err != nil {
		return ""
	}

	// Resolve the link using the reference URL
	resolved := refURL.ResolveReference(url)

	// Strip fragments as we only care about unique pages
	resolved.Fragment = ""

	// Strip trailing slashes as we only care about unique pages
	return strings.TrimSuffix(resolved.String(), "/")
}
