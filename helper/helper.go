package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse sends a JSON response with an error message and status code.
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// WriteJSON writes a JSON response with the given status code and content.
func WriteJSON(respw http.ResponseWriter, statusCode int, content any) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)

	// Convert content to JSON and handle potential errors
	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(respw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respw.Write(jsonData)
}
