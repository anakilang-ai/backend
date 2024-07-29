package modules

import (
	"net/http"
)

// List of allowed origins
var Origins = []string{
	"http://localhost:8080",
	"https://anakilang-ai.github.io/",
}

// Checks if the origin is allowed
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}

// Sets CORS headers for the response
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if !isAllowedOrigin(origin) {
		// If the origin is not allowed, do not set CORS headers
		return false
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Login")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return true
	}

	return false
}
