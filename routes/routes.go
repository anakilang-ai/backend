package routes

import (
	"net/http"

	"github.com/anakilang-ai/backend/modules"
	controller "github.com/anakilang-ai/backend/controller"
	"github.com/anakilang-ai/backend/helper"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if modules.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}

	if modules.ErrorMongoconn != nil {
		helper.ErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : database, " + modules.ErrorMongoconn.Error())
		return
	}

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		Home(w, r)
	case method == "POST" && path == "/signup":
		controller.SignUp(modules.Mongoconn, "users", w, r)
	case method == "POST" && path == "/login":
		controller.LogIn(modules.Mongoconn, w, r, modules.GetEnv("PASETOPRIVATEKEY"))
	case method == "POST" && path == "/chat":
		controller.Chat(w, r, modules.GetEnv("TOKENMODEL"))
	default:
		helper.ErrorResponse(w, r, http.StatusNotFound, "Not Found", "The requested resource was not found")
	}
}

func Home(respw http.ResponseWriter, req *http.Request) {
	resp := map[string]string{
		"github_repo": "https://github.com/anakilang-ai/backend",
		"message": "",
	}
	helper.WriteJSON(respw, http.StatusOK, resp)
}