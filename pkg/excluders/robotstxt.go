package excluders

import (
	"fmt"
	"net/http"
	urlpkg "net/url"

	"github.com/temoto/robotstxt"
)

// RobotsTxtExcluder bases its exclusion rules on a given robots.txt file
type RobotsTxtExcluder struct {
	enable bool
	rules  *robotstxt.Group
}

func NewRobotsTxtExcluder(agent string, enable bool, url *urlpkg.URL) (*RobotsTxtExcluder, error) {
	group, err := load(agent, url)
	if err != nil {
		return nil, fmt.Errorf("could not load: %w", err)
	}

	return &RobotsTxtExcluder{
		enable: enable,
		rules:  group,
	}, nil
}

// Exclude checks against the rules found in the robots.txt file to see if a given path
// should be excluded. It allows all paths if disabled.
func (t *RobotsTxtExcluder) Exclude(path string) bool {
	if !t.enable {
		return false
	}

	// Test returns true for allowed paths, so we negate it
	return !t.rules.Test(path)
}

// load fetches the URL's robots.txt file and parses its rules for the given agent. If no
// robots.txt file can be found, it can be assumed that no rules will be loaded.
func load(agent string, url *urlpkg.URL) (*robotstxt.Group, error) {
	rtPath, err := urlpkg.Parse("/robots.txt")
	if err != nil {
		return nil, err
	}
	rtURL := url.ResolveReference(rtPath).String()

	resp, err := http.Get(rtURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	return data.FindGroup(agent), nil
}
