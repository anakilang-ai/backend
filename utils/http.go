package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse menangani respons error JSON dengan pesan yang sesuai.
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	writeJSON(respw, statusCode, resp)
}

// WriteJSON menulis response JSON ke ResponseWriter dengan status code yang diberikan.
func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	err := json.NewEncoder(respw).Encode(content)
	if err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

// jsonStr mengembalikan representasi JSON dari objek sebagai string.
func jsonStr(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return ""
	}
	return string(jsonData)
}
