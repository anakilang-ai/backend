package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/anakilang-ai/backend/routes"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		routes.URL(w, r)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		routes.URL(w, r)
	})
	r.Put("/", func(w http.ResponseWriter, r *http.Request) {
		routes.URL(w, r)
	})
	r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		routes.URL(w, r)
	})

	port := ":8080"
	http.ListenAndServe(port, r)
}