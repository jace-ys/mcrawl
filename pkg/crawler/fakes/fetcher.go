package fakes

import "fmt"

type FakeFetcher struct {
	results map[string][]string
}

func NewFakeFetcher() *FakeFetcher {
	return &FakeFetcher{
		results: map[string][]string{
			"https://golang.org": []string{
				"https://golang.org/pkg",
				"https://golang.org/cmd",
				"https://monzo.com",
			},
			"https://golang.org/pkg": []string{
				"https://golang.org",
				"https://golang.org/cmd",
				"https://golang.org/pkg/fmt",
				"https://golang.org/pkg/os",
			},
			"https://golang.org/pkg/fmt": []string{
				"https://golang.org",
				"https://golang.org/pkg",
			},
			"https://golang.org/pkg/os": []string{
				"https://golang.org",
				"https://golang.org/pkg",
			},
			"https://golang.org/cmd": []string{},
		},
	}
}

func (f *FakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f.results[url]; ok {
		return res, nil
	}

	return nil, fmt.Errorf("not found: %s", url)
}
