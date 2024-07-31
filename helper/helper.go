package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse mengirimkan respons JSON berisi pesan kesalahan
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// WriteJSON mengirimkan respons JSON dengan konten yang diberikan
func WriteJSON(respw http.ResponseWriter, statusCode int, content any) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

// Jsonstr mengubah struktur data menjadi string JSON
func Jsonstr(strc any) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}
