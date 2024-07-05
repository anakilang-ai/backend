package ailang

import (
	"github.com/anakilang-ai/backend/routes"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WebHook", routes.URL)
}