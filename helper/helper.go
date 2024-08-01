package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse sends a JSON error response with the specified status code, error, and message.
func ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err, msg string) {
	response := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(w, statusCode, response)
}

// WriteJSON sends a JSON response with the specified status code and content.
// It handles errors in JSON marshalling gracefully and logs the error.
func WriteJSON(w http.ResponseWriter, statusCode int, content any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, err := json.Marshal(content)
	if err != nil {
		// Log the error and send a generic error response
		log.Printf("Failed to marshal JSON response: %v", err)
		http.Error(w, `{"error":"Internal Server Error","message":"Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// JSONString converts a Go struct or value to its JSON string representation.
// It handles errors gracefully by logging and returning an empty string.
func JSONString(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Log the error and return an empty string
		log.Printf("Failed to marshal JSON: %v", err)
		return ""
	}
	return string(jsonData)
}
