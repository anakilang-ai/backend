package routes

import (
	"net/http"

	"github.com/anakilang-ai/backend/controller"
	"github.com/anakilang-ai/backend/helper"
	"github.com/anakilang-ai/backend/modules"
)

// URL is the main routing handler for all requests.
func URL(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers and return early if it's a preflight request.
	if modules.SetAccessControlHeaders(w, r) {
		return
	}

	// Check for MongoDB connection error.
	if modules.ErrorMongoconn != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "database connection error: "+modules.ErrorMongoconn.Error())
		return
	}

	// Route based on HTTP method and URL path.
	method, path := r.Method, r.URL.Path
	switch {
	case method == http.MethodGet && path == "/":
		Home(w, r)
	case method == http.MethodPost && path == "/signup":
		controller.SignUp(modules.Mongoconn, "users", w, r)
	case method == http.MethodPost && path == "/login":
		controller.LogIn(modules.Mongoconn, w, r, modules.GetEnv("PASETOPRIVATEKEY"))
	case method == http.MethodPost && path == "/chat":
		controller.Chat(w, r, modules.GetEnv("TOKENMODEL"))
	default:
		helper.ErrorResponse(w, r, http.StatusNotFound, "Not Found", "The requested resource was not found")
	}
}

// Home handles the root route.
func Home(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"github_repo": "https://github.com/anakilang-ai/backend",
		"message":     "Welcome to the API",
	}
	helper.WriteJSON(w, http.StatusOK, resp)
}
