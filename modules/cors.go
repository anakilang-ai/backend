package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Fungsi untuk memuat variabel lingkungan dari file .env
func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}

// Fungsi untuk mengambil nilai dari variabel lingkungan yang diberi nama oleh kunci
// Mengembalikan string kosong jika variabel tersebut tidak ada
func GetEnv(envName string) string {
	return os.Getenv(envName)
}

// Daftar origins yang diizinkan
var Origins = []string{
	"http://localhost:8080",
	"https://anakilang-ai.github.io/",
}

// Fungsi untuk memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}

// Fungsi untuk mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	if isAllowedOrigin(origin) {
		// Set CORS headers for the preflight request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Login")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return true
		}
		// Set CORS headers for the main request.
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	return false
}

// Contoh handler untuk HTTP server
func handler(w http.ResponseWriter, r *http.Request) {
	if SetAccessControlHeaders(w, r) {
		return
	}

	// Lanjutkan penanganan request jika bukan preflight request
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	// Memuat variabel lingkungan dari file .env
	LoadEnv()

	// Mengambil dan mencetak variabel lingkungan sebagai contoh
	value := GetEnv("YOUR_ENV_VARIABLE")
	fmt.Println("Nilai dari YOUR_ENV_VARIABLE:", value)

	// Menjalankan HTTP server
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
