// config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file if it exists.
// It logs an error message if the .env file cannot be found or loaded.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file does not exist")
		} else {
			log.Printf("Error loading .env file: %v", err)
		}
	}
}

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the default value if the key does not exist or is empty.
// If no default value is provided and the key does not exist, it returns an empty string.
func GetEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
