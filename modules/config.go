package modules

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(envName string) string {
	envFile, _ := godotenv.Read("../.env")
	return envFile[envName]
	return os.Getenv(envName)
}
