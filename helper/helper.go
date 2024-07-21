package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}
	respw.Write(jsonData)
}

func jsonStr(strc interface{}) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return ""
	}
	return string(jsonData)
}
