package modules

//mengimpor package os dari standard library Go.
import (
	"os"
	// "github.com/joho/godotenv"
)

// mengambil nilai dari variabel lingkungan dengan nama envName.
func GetEnv(envName string) string {
	// envFile, _ := godotenv.Read("../.env")
	// return envFile[envName]
	return os.Getenv(envName)
}
