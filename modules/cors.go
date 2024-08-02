package modules

import (
	"net/http"
)

// Allowed origins for CORS requests.
var AllowedOrigins = map[string]struct{}{
	"http://localhost:8080":          {},
	"http://127.0.0.1:5501":          {},
	"https://anakilang-ai.github.io": {},
}

// Checks if the provided origin is allowed.
func isAllowedOrigin(origin string) bool {
	_, allowed := AllowedOrigins[origin]
	return allowed
}

// Sets CORS headers for the response.
// Returns true if the request was handled as a preflight request.
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		if r.Method == http.MethodOptions {
			// Handle preflight request.
			setCORSHeaders(w, origin, true)
			w.WriteHeader(http.StatusNoContent)
			return true
		}
		// Handle main request.
		setCORSHeaders(w, origin, false)
		return false
	}

	// If the origin is not allowed, no CORS headers are set.
	return false
}

// Sets CORS headers for the response.
func setCORSHeaders(w http.ResponseWriter, origin string, isPreflight bool) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
	if isPreflight {
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
}
