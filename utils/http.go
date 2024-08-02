package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse formats and sends an error response in JSON format.
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// WriteJSON sends a JSON response with the specified status code.
func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	if err := json.NewEncoder(respw).Encode(content); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Jsonstr converts an interface to a JSON string. Logs and returns an empty string if there is an error.
func Jsonstr(strc interface{}) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return ""
	}
	return string(jsonData)
}
