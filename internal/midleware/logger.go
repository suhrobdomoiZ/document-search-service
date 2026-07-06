package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter

	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("Started to perform request", "method", request.Method, "path", request.URL.Path)

		start := time.Now()
		wrapper := &responseWriter{ResponseWriter: writer, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, request)

		duration := time.Since(start)

		slog.Info("request completed",
			"method", request.Method,
			"path", request.URL.Path,
			"status", wrapper.statusCode,
			"duration", duration,
		)
	})
}
