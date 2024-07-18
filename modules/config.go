package modules

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from the .env file if it exists
func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: .env file not found or could not be read, proceeding with environment variables from the OS: %v", err)
	}
}

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
func GetEnv(envName string) string {
	return os.Getenv(envName)
}
