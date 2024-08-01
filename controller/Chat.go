package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/anakilang-ai/backend/helper"
	"github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/modules"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Logger instance
var log = logrus.New()

func init() {
	// Customize the logger if needed
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

// Chat handles AI chat requests and interacts with the Hugging Face API
func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
	var chat models.AIRequest

	// Decode the request body into the chat struct
	if err := json.NewDecoder(req.Body).Decode(&chat); err != nil {
		log.WithError(err).Error("Failed to parse request body")
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
		return
	}

	// Validate the chat prompt
	if chat.Prompt == "" {
		log.Warn("Empty prompt in request")
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Prompt cannot be empty")
		return
	}

	client := resty.New()

	// Hugging Face API URL and token
	apiUrl := modules.GetEnv("HUGGINGFACE_API_URL", "https://api.huggingface.co/v1/models")
	apiToken := "Bearer " + tokenmodel

	// Request to Hugging Face API with retry mechanism
	response, err := makeRequestWithRetry(client, apiUrl, apiToken, chat.Prompt)
	if err != nil {
		log.WithError(err).Error("Failed to get a valid response from Hugging Face API")
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	generatedText, err := extractGeneratedText(response.Body())
	if err != nil {
		log.WithError(err).Error("Failed to extract generated text")
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"answer": generatedText})
}

// makeRequestWithRetry performs a request with retry logic in case of model loading
func makeRequestWithRetry(client *resty.Client, url, token, prompt string) (*resty.Response, error) {
	var response *resty.Response
	var err error
	maxRetries := 5
	retryDelay := 20 * time.Second

	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		response, err = makeRequest(client, url, token, prompt)
		if err != nil {
			log.WithError(err).Error("Request failed")
		} else if response.StatusCode() == http.StatusOK {
			return response, nil
		} else if isModelLoading(response.Body()) {
			log.Info("Model is currently loading, retrying...")
			time.Sleep(retryDelay)
		} else {
			log.WithField("status_code", response.StatusCode()).Error("Error from Hugging Face API")
			return response, errors.New("error from Hugging Face API: " + string(response.Body()))
		}
	}

	return nil, errors.New("maximum retries reached")
}

// makeRequest performs a single API request
func makeRequest(client *resty.Client, url, token, prompt string) (*resty.Response, error) {
	return client.R().
		SetHeader("Authorization", token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"inputs": prompt}).
		Post(url)
}

// isModelLoading checks if the response indicates that the model is loading
func isModelLoading(body []byte) bool {
	var errorResponse map[string]interface{}
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorMessage, ok := errorResponse["error"].(string); ok && errorMessage == "Model is currently loading" {
			return true
		}
	}
	return false
}

// extractGeneratedText extracts the generated text from the API response
func extractGeneratedText(body []byte) (string, error) {
	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if len(data) > 0 {
		if generatedText, ok := data[0]["generated_text"].(string); ok {
			return generatedText, nil
		}
	}
	return "", errors.New("error extracting generated text")
}
