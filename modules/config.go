package modules

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadEnv loads environment variables from a .env file if it exists.
// It returns an error if it fails to load the file or parse it.
func LoadEnv(filePath string) error {
	err := godotenv.Load(filePath)
	if err != nil {
		return err
	}
	return nil
}

// GetEnv retrieves the value of the environment variable named by envName.
// It returns the value of the environment variable or an empty string if it is not set.
func GetEnv(envName string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		log.Printf("Warning: Environment variable %s not set", envName)
	}
	return value
}
