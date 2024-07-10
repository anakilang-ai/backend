package main

//mengimpor beberapa library yang dibutuhkan untuk menjalankan fungsi-fungsi yang ada di file tersebut.
import (
	"fmt"
	"net/http"

	"github.com/anakilang-ai/backend/routes"
)

// mendaftarkan fungsi routes.URL untuk menangani semua request yang masuk ke path "/" (root path).  Kemudian program dijalankan pada port 8080 dan menampilkan pesan bahwa server telah aktif di http://localhost:8080.
func main() {
	http.HandleFunc("/", routes.URL)
	port := ":8080"
	fmt.Println("Server started at: http://localhost" + port)
	http.ListenAndServe(port, nil)
}
