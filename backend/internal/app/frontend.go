package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	urlpath "path"
	"path/filepath"
	"strings"
)

func frontendFallback(next http.Handler, distDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)
		if recorder.Code != http.StatusNotFound {
			copyHeaders(w.Header(), recorder.Header())
			w.WriteHeader(recorder.Code)
			_, _ = io.Copy(w, recorder.Body)
			return
		}
		if shouldSkipFrontendFallback(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		serveFrontendFile(w, r, distDir)
	})
}

func shouldSkipFrontendFallback(path string) bool {
	return strings.HasPrefix(path, "/api/") ||
		path == "/api" ||
		strings.HasPrefix(path, "/schemas/") ||
		path == "/openapi.json" ||
		strings.HasPrefix(path, "/docs")
}

func serveFrontendFile(w http.ResponseWriter, r *http.Request, distDir string) {
	if distDir == "" {
		http.NotFound(w, r)
		return
	}
	cleanPath := urlpath.Clean("/" + r.URL.Path)
	target := filepath.Join(distDir, filepath.FromSlash(strings.TrimPrefix(cleanPath, "/")))
	info, err := os.Stat(target)
	if err == nil && !info.IsDir() {
		http.ServeFile(w, r, target)
		return
	}
	indexPath := filepath.Join(distDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, indexPath)
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
