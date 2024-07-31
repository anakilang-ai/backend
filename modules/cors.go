package modules

import (
	"net/http"
)

// Daftar origins yang diizinkan
var AllowedOrigins = []string{
	"http://localhost:8080",
	"https://anakilang-ai.github.io/",
}

// Fungsi untuk memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, allowed := range AllowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

// Fungsi untuk mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Periksa jika origin diizinkan
	if isAllowedOrigin(origin) {
		// Set header untuk permintaan preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") // Mengganti 'Login' dengan 'Authorization' untuk kepatuhan standar
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT,OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return true
		}
		// Set header untuk permintaan utama
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		return false
	}

	return false
}
