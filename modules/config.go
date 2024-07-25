// config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file if it exists.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
}

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the default value if the key does not exist.
func GetEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
