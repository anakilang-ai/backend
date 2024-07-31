package routes

import (
	"fmt"
	"net/http"
)

// URL adalah handler untuk rute utama
func URL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Selamat datang di server Go!")
}
