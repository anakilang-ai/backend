package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anakilang-ai/backend/helper"
	"github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/modules"
	"github.com/go-resty/resty/v2"
)

func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
	var chat models.AIRequest

	// Decode the request body into the chat struct
	if err := json.NewDecoder(req.Body).Decode(&chat); err != nil {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body: "+err.Error())
		return
	}

	// Validate the chat prompt
	if chat.Prompt == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	client := resty.New()

	// Hugging Face API URL and token
	apiUrl := modules.GetEnv("HUGGINGFACE_API_URL")
	apiToken := "Bearer " + tokenmodel

	var response *resty.Response
	var err error
	maxRetries := 5
	retryDelay := 20 * time.Second

	// Request to Hugging Face API with retry mechanism
	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		response, err = client.R().
			SetHeader("Authorization", apiToken).
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{"inputs": chat.Prompt}).
			Post(apiUrl)

		if err != nil {
			log.Printf("Error making request: %v", err)
			helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error making request: "+err.Error())
			return
		}

		if response.StatusCode() == http.StatusOK {
			break
		}

		// Check if the error is due to model loading
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(response.Body(), &errorResponse); err == nil {
			if errorMessage, ok := errorResponse["error"].(string); ok && errorMessage == "Model is currently loading" {
				time.Sleep(retryDelay)
				continue
			}
		}

		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error from Hugging Face API: "+string(response.Body()))
		return
	}

	if response.StatusCode() != http.StatusOK {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error from Hugging Face API: "+string(response.Body()))
		return
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error parsing response body: "+err.Error())
		return
	}

	if len(data) > 0 {
		generatedText, ok := data[0]["generated_text"].(string)
		if !ok {
			helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error extracting generated text")
			return
		}
		helper.WriteJSON(respw, http.StatusOK, map[string]string{"answer": generatedText})
	} else {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: response")
	}
}
