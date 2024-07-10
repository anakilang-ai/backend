package ailang

// memuat rute dari paket routes dan fungsi dari paket functions dalam sebuah proyek Go.
import (
	"github.com/anakilang-ai/backend/routes"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// menginisialisasi dan mendaftarkan handler HTTP bernama "WebHook" yang akan mengeksekusi fungsi routes.URL pada saat aplikasi dijalankan.
func init() {
	functions.HTTP("WebHook", routes.URL)
}
