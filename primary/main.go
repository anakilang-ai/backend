package main

//mengimpor beberapa library yang dibutuhkan untuk menjalankan fungsi-fungsi yang ada di file tersebut.
import (
	"fmt"
	"net/http"

	"github.com/anakilang-ai/backend/routes"
)

func main() {
	http.HandleFunc("/", routes.URL)
	port := ":8080"
	fmt.Println("Server started at: http://localhost" + port)
	http.ListenAndServe(port, nil)
}
