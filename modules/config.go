package modules

import (
	"log"
	"os"
)

// GetEnv retrieves the value of the environment variable named by envName.
// It returns an empty string if the environment variable is not set.
func GetEnv(envName string) string {
	// Mengambil nilai variabel lingkungan
	value := os.Getenv(envName)
	if value == "" {
		// Mencetak peringatan jika variabel lingkungan tidak diatur
		log.Printf("Warning: Environment variable %s is not set.", envName)
	}
	return value
}
