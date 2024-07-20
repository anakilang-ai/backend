package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse mengirimkan respon error dalam format JSON dengan status code yang sesuai
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// WriteJSON menulis respon dalam format JSON dengan status code yang diberikan
func WriteJSON(respw http.ResponseWriter, statusCode int, content any) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	jsonData, err := json.Marshal(content)
	if err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}
	respw.Write(jsonData)
}

// Jsonstr mengubah struktur data ke dalam format string JSON
func Jsonstr(strc any) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return ""
	}
	return string(jsonData)
}
