package routes

import (
	"fmt"
	"net/http"
)

// URL adalah handler untuk rute utama ("/")
func URL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Selamat datang di server!")
	} else {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
	}
}
