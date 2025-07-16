package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"
)

// LoggingResponseWriter wraps http.ResponseWriter to capture status code and bytes written.
type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode   int
	BytesWritten int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.BytesWritten += n
	return n, err
}

// Start Initialize and run proxy servers
func Start(listenAddr string) error {
	proxy := &httputil.ReverseProxy{
		// modifying requests with Director
		Director: func(req *http.Request) {
			// record request start time in header
			req.Header.Set("X-Start-Time", time.Now().UTC().Format(time.RFC3339Nano))
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lrw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		destination := r.URL
		if !destination.IsAbs() {
			http.Error(lrw, "Absolute URL required", http.StatusBadRequest)
			return
		}
		r.Host = destination.Host

		// lrw include http.ResponseWriter
		proxy.ServeHTTP(lrw, r)

		duration := time.Since(startTime)
		// memory check
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// structured JSON style log
		log.Printf(
			`{"method":"%s","url":"%s","startTime":"%s","durationMs":%d,"status":%d,"bytes":%d,"memAlloc":%d}`,
			r.Method,
			r.URL.String(),
			startTime.UTC().Format(time.RFC3339Nano),
			duration.Milliseconds(),
			lrw.StatusCode,
			lrw.BytesWritten,
			m.Alloc,
		)
	})

	log.Printf("Starting HTTP proxy at %s", listenAddr)
	return http.ListenAndServe(listenAddr, handler)
}
