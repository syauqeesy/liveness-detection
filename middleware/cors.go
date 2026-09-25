package middleware

import (
	"net/http"

	"github.com/syauqeesy/liveness-detection/configuration"
)

func Cors(configuration *configuration.Configuration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Origin") == configuration.Application.Client {
				w.Header().Set("Access-Control-Allow-Origin", configuration.Application.Client)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			next.ServeHTTP(w, r)
		})
	}
}
