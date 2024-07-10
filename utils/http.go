package utils

// memuat paket encoding/json untuk encoding dan decoding JSON, paket log untuk logging, dan paket net/http untuk membuat aplikasi HTTP dalam proyek Go.
import (
	"encoding/json"
	"log"
	"net/http"
)

// mengirimkan respons HTTP dengan status kode tertentu dan pesan kesalahan dalam format JSON.
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// mengirimkan respons HTTP dengan status kode tertentu dan konten yang telah dikonversi menjadi JSON.
func WriteJSON(respw http.ResponseWriter, statusCode int, content any) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

func Jsonstr(strc any) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}
