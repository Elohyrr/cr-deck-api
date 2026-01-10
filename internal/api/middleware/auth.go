package middleware

import (
	"net/http"
	"strings"
)

// Auth validates Bearer token authentication
func Auth(apiToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"error": "unauthorized", "message": "missing Bearer token"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token != apiToken {
				http.Error(w, `{"error": "unauthorized", "message": "invalid token"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
