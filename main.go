package ailang

// memuat rute dari paket routes dan fungsi dari paket functions dalam sebuah proyek Go.
import (
	"github.com/anakilang-ai/backend/routes"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WebHook", routes.URL)
}
