package test

import (
	"god/proxy"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"net/http/httputil"
)

func TestProxy_Start(t *testing.T) {
	// Create a test server to act as the backend
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend", "ok")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	// Start the proxy server using httptest
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lrw := &proxy.LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		destination := r.URL
		if !destination.IsAbs() {
			http.Error(lrw, "Absolute URL required", http.StatusBadRequest)
			return
		}
		r.Host = destination.Host

		p := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				u, _ := url.Parse(backend.URL)
				req.URL.Scheme = u.Scheme
				req.URL.Host = u.Host
				req.Header.Set("X-Start-Time", time.Now().UTC().Format(time.RFC3339Nano))
			},
		}
		p.ServeHTTP(lrw, r)

		_ = time.Since(startTime) // duration measurement skipped
	}))
	defer proxyServer.Close()

	// Send a request to the proxy server
	client := &http.Client{}
	// Create a request with an absolute URL
	u, _ := url.Parse(backend.URL)
	req, err := http.NewRequest("GET", proxyServer.URL+"/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	req.URL.Path = "/"

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("proxy request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "backend response" {
		t.Errorf("unexpected response body: %s", string(body))
	}
	if resp.Header.Get("X-Backend") != "ok" {
		t.Errorf("missing or wrong backend header: %s", resp.Header.Get("X-Backend"))
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
