package modules

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file if it exists.
// This function is optional and can be used for local development.
func LoadEnv(filePath string) error {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		if err := godotenv.Load(filePath); err != nil {
			return err
		}
	}
	return nil
}

// GetEnv retrieves the value of the environment variable named by envName.
// It returns an empty string if the variable is not set.
func GetEnv(envName string) string {
	return os.Getenv(envName)
}

// GetEnvOrDefault retrieves the value of the environment variable named by envName.
// If the variable is not set, it returns the provided default value.
func GetEnvOrDefault(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	return value
}

// Example usage:
// func main() {
//     err := LoadEnv(".env")
//     if err != nil {
//         log.Fatalf("Error loading .env file: %v", err)
//     }
//
//     port := GetEnvOrDefault("PORT", "8080")
//     log.Printf("Port: %s", port)
// }
