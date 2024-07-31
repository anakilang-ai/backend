package modules

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Inisialisasi untuk membaca file .env
func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Kesalahan saat memuat file .env: %v", err)
	}
}

// Fungsi untuk mendapatkan nilai variabel lingkungan
func GetEnv(envName string) string {
	return os.Getenv(envName)
}
