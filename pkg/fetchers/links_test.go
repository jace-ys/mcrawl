package fetchers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchHTML(t *testing.T) {
	var (
		contentType string
		statusCode  int
		body        string
	)

	fetcher := LinksFetcher{}

	handler := http.NewServeMux()
	server := httptest.NewServer(handler)
	defer server.Close()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)

		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	})

	tt := []struct {
		name        string
		contentType string
		statusCode  int
		body        string
		err         error
	}{
		{
			name:        "Success",
			contentType: "text/html",
			statusCode:  http.StatusOK,
			body:        "Gopher",
		},
		{
			name:        "Invalid Response Code",
			contentType: "text/html",
			statusCode:  http.StatusNotFound,
			err:         ErrInvalidResponseCode,
		},
		{
			name:        "Invalid Content-Type",
			contentType: "application/json",
			statusCode:  http.StatusOK,
			err:         ErrInvalidContentType,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			contentType = tc.contentType
			statusCode = tc.statusCode
			body = tc.body

			data, err := fetcher.fetchHTML(server.URL)
			if tc.err != nil {
				assert.True(t, errors.Is(err, tc.err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.body, string(data))
			}
		})
	}
}

func TestFindLinks(t *testing.T) {
	fetcher := LinksFetcher{}

	data := `
<html>
	<a href="/"></a>
	<a href="https://monzo.com"></a>
	<a href="overdrafts"></a>
	<a href="../about"></a>
	<a href="/features/travel/"></a>
	<a href="#home"></a>
	<a href="/blog?page=2"></a>
	<a href="//monzo.com/careers"></a>
	<a></a>
</html>`

	found, err := fetcher.findLinks([]byte(data))

	assert.NoError(t, err)
	assert.Len(t, found, 8)
}

func TestSanitise(t *testing.T) {
	fetcher := LinksFetcher{}
	baseURL, err := url.Parse("https://monzo.com/i/loans")
	assert.NoError(t, err)

	tt := []struct {
		name     string
		link     string
		expected string
	}{
		{
			name:     "Root Path",
			link:     "/",
			expected: "https://monzo.com",
		},
		{
			name:     "Absolute URL",
			link:     "https://monzo.com",
			expected: "https://monzo.com",
		},
		{
			name:     "Relative Path",
			link:     "overdrafts",
			expected: "https://monzo.com/i/overdrafts",
		},
		{
			name:     "Relative Path Parent",
			link:     "../about",
			expected: "https://monzo.com/about",
		},
		{
			name:     "Trailing Slash",
			link:     "/features/travel/",
			expected: "https://monzo.com/features/travel",
		},
		{
			name:     "Hash Link",
			link:     "#home",
			expected: "https://monzo.com/i/loans",
		},
		{
			name:     "Query Parameters",
			link:     "/blog?page=2",
			expected: "https://monzo.com/blog?page=2",
		},
		{
			name:     "Missing Scheme",
			link:     "//monzo.com/careers",
			expected: "https://monzo.com/careers",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := fetcher.sanitise(tc.link, baseURL)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
