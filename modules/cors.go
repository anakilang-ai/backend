package modules

import (
	"net/http"
)

// Daftar origins yang diizinkan
var allowedOrigins = []string{
	"http://localhost:8080",
	"https://anakilang-ai.github.io/",
}

// isAllowedOrigin memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

// SetAccessControlHeaders mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

// Middleware CORS untuk HTTP handler
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetAccessControlHeaders(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
