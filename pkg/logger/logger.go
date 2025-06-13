package logger

import (
	"fmt"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func ApiInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		sw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		duration := time.Since(startTime).Milliseconds()
		statusText := http.StatusText(sw.status)

		currTime := time.Now().Format("2006-01-02 15:04:05")

		fmt.Printf("[%s] %s %s - %d %s (%dms)\n", currTime, r.Method, r.RequestURI, sw.status, statusText, duration)
	})
}
