package modules

import (
	"net/http"
)

// Daftar asal yang diizinkan
var allowedOrigins = []string{
	"http://localhost:8080",
	"https://anakilang-ai.github.io/",
}

// isAllowedOrigin memeriksa apakah asal yang diberikan ada dalam daftar asal yang diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

// SetAccessControlHeaders mengatur header CORS dan menangani permintaan preflight
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		// Mengatur header CORS untuk semua permintaan
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)

		// Menangani permintaan preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return true
		}
	}

	return false
}
