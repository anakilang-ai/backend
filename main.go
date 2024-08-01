package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/anakilang-ai/backend/routes"
)

func main() {
	r := mux.NewRouter()

	// Menggunakan handler untuk semua metode HTTP
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.URL(w, r)
	})

	port := ":8080"
	http.ListenAndServe(port, r)
}