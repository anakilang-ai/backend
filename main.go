package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/anakilang-ai/backend/routes"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Menggunakan handler untuk semua metode HTTP
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		routes.URL(w, r)
	})

	port := ":8080"
	http.ListenAndServe(port, r)
}
