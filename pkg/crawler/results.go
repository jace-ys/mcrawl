package crawler

import "sync"

// Result is a struct used for storing the links found for each target URL crawled
type Result struct {
	Target string
	Links  []string
	Err    error
}

// resultsMap is a synchronised map used for caching the links found for each URL crawled
type resultsMap struct {
	urls map[string]Result
	mu   sync.RWMutex
}

func makeResultsMap() *resultsMap {
	return &resultsMap{
		urls: make(map[string]Result),
	}
}

func (rm *resultsMap) add(r Result) {
	rm.mu.Lock()
	rm.urls[r.Target] = r
	rm.mu.Unlock()
}

func (rm *resultsMap) contains(url string) bool {
	rm.mu.RLock()
	_, ok := rm.urls[url]
	rm.mu.RUnlock()

	return ok
}
