package modul

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv memuat variabel lingkungan dari file .env.
func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}

// GetEnv mengambil nilai dari variabel lingkungan yang diberi nama oleh kunci.
// Mengembalikan string kosong jika variabel tersebut tidak ada.
func GetEnv(envName string) string {
	return os.Getenv(envName)
}

package main

import (
	"fmt"
	"jalur_proyek_anda/modul"
)

func main() {
	// Memuat variabel lingkungan dari file .env
	modul.LoadEnv()

	// Mengambil variabel lingkungan
	nilai := modul.GetEnv("YOUR_ENV_VARIABLE")
	fmt.Println("Nilai dari YOUR_ENV_VARIABLE:", nilai)
}
