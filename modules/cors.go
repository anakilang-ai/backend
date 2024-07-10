package modules

//mengimpor package net/http dari standard library Go.  Package net/http menyediakan berbagai fungsi untuk membangun dan menjalankan server HTTP, menangani request dan response HTTP, dan melakukan crawling web.
import (
	"net/http"
)

// Mendaftar origins yang diizinkan
var Origins = []string{
	"http://localhost:8080",
	"https://anakilang-ai.github.io/",
}

// mengecek apakah sebuah origin (biasanya URL asal request) diizinkan untuk mengakses resource.
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}

// mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		// Tetapkan header CORS untuk permintaan preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return true
		}

		// Tetapkan header CORS untuk permintaan utama.
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		return false
	}

	return false
}
