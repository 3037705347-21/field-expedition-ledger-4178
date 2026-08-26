package httpapi

import (
	"net/http"
	"os"
	"strings"
)

func staticHandler() http.Handler {
	root := "web"
	if _, err := os.Stat(root); err != nil {
		root = "."
	}
	fileServer := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/":
			r.URL.Path = "/index.html"
		case strings.HasPrefix(r.URL.Path, "/static/"):
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/static")
		default:
			http.NotFound(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
