package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse sends a JSON response with an error message.
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// WriteJSON writes a JSON response with the specified status code.
func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)

	// Marshal the content to JSON
	jsonData, err := json.Marshal(content)
	if err != nil {
		// If marshaling fails, log the error and return a generic error response
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(respw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response writer
	_, err = respw.Write(jsonData)
	if err != nil {
		// If writing fails, log the error
		log.Printf("Error writing response: %v", err)
	}
}

// Jsonstr converts an interface{} to a JSON string.
func Jsonstr(strc interface{}) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		// If marshaling fails, log the error and return an empty string
		log.Printf("Error marshalling to JSON: %v", err)
		return ""
	}
	return string(jsonData)
}
