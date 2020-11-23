package e2e

import (
	"net/http"
	"net/http/httptest"
)

func startMockServer() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		body := `
# robotstxt.org/

User-agent: *
Disallow: /forbidden`

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(body))
	})

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := `
<html>
  <h1>Home</h1>
  <a href="/about">About</a>
  <a href="/forbidden">Forbidden</a>
  <a href="https://monzo.com">First</a>
</html>`

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(body))
	})

	handler.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		body := `
<html>
	<h1>About</h1>
	<a href="/">Index</a>
	<a href="/contact">Contact</a>
</html>`

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(body))
	})

	handler.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		body := `
<html>
	<h1>Contact</h1>
</html>`

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(body))
	})

	return httptest.NewServer(handler)
}
