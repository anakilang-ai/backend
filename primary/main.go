package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anakilang-ai/backend/routes"
)

func main() {
	// Menangani route utama
	http.HandleFunc("/", routes.URL)

	// Menentukan port server
	port := ":8080"
	fmt.Printf("Server started at: http://localhost%s\n", port)

	// Menjalankan server dan menangani kesalahan jika terjadi
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server gagal untuk memulai: %v", err)
	}
}
