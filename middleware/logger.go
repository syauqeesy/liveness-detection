package middleware

import (
	"net/http"
	"time"
	"uuid"

	"github.com/liveness-detection/common"
)

func Logger(logger common.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			id := r.Header.Get("X-Request-Id")
			if id == "" {
				id = uuid.NewV7().String()
			}

			w.Header().Set("X-Request-Id", id)

			logger.Info(
				"http request started",
				"request_id", id,
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(w, r)

			logger.Info(
				"http request completed",
				"request_id", id,
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
