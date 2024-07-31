package modules

import (
	"os"
	"log"
)

// GetEnv retrieves the value of the environment variable named by envName.
// It returns an empty string if the environment variable is not set.
func GetEnv(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Printf("Warning: Environment variable %s is not set.", envName)
	}
	return value
}
