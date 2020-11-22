package excluders

import (
	urlpkg "net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/temoto/robotstxt"
)

var (
	startURL, _ = urlpkg.Parse("https://golang.org/")
)

func TestExclude(t *testing.T) {
	testRobotsTxt := `
# robotstxt.org/

User-agent: *
Disallow: /login
Disallow: /users/`

	data, err := robotstxt.FromString(testRobotsTxt)
	assert.NoError(t, err)

	tt := []struct {
		name     string
		enable   bool
		path     string
		expected bool
	}{
		{
			name:     "Exact Match",
			enable:   true,
			path:     "/login",
			expected: true,
		},
		{
			name:     "Partial Match",
			enable:   true,
			path:     "/users/gopher",
			expected: true,
		},
		{
			name:     "No Match",
			enable:   true,
			path:     "/home",
			expected: false,
		},
		{
			name:     "Disable",
			enable:   false,
			path:     "/login",
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			excluder := RobotsTxtExcluder{
				enable: tc.enable,
				rules:  data.FindGroup("MCrawl"),
			}

			exclude := excluder.Exclude(tc.path)
			assert.Equal(t, tc.expected, exclude)
		})
	}
}
