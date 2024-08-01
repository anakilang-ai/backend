package modules

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
// It returns an error if the file cannot be loaded or if there's an issue reading it.
func LoadEnv(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		if err := godotenv.Load(filePath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// GetEnv retrieves the value of the environment variable named by envName.
// If the environment variable is not set, it returns a default value if provided.
func GetEnv(envName string, defaultValue ...string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		log.Printf("Warning: Environment variable %s not set", envName)
		return ""
	}
	return value
}
