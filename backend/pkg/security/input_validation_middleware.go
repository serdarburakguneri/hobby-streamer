package security

import (
	"net/http"
	"strings"
)

func InputValidationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > 10*1024*1024 { // 10 MB max
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				_, _ = w.Write([]byte(`{"error": "Request too large"}`))
				return
			}

			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				ct := r.Header.Get("Content-Type")
				if !strings.Contains(ct, "application/json") &&
					!strings.Contains(ct, "multipart/form-data") &&
					!strings.Contains(ct, "application/x-www-form-urlencoded") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnsupportedMediaType)
					_, _ = w.Write([]byte(`{"error": "Unsupported content type"}`))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
