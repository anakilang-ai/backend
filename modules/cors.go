package modules

import (
	"net/http"
)

// AllowedOrigins lists the origins that are permitted to access the resources.
var AllowedOrigins = []string{
	"http://localhost:8080",
	"http://127.0.0.1:5501",
	"https://anakilang-ai.github.io",
}

// isAllowedOrigin checks if the given origin is in the list of allowed origins.
func isAllowedOrigin(origin string) bool {
	for _, allowedOrigin := range AllowedOrigins {
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}

// SetAccessControlHeaders sets CORS headers based on the request's origin.
// It returns true if the request was handled (preflight request) and false otherwise.
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return true
		}
		// Handle actual requests
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Origin not allowed
	return false
}
