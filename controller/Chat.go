package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/modules"
	"github.com/anakilang-ai/backend/utils"
	"github.com/go-resty/resty/v2"
)

// LoadDatasetLocal loads the dataset from a local CSV file and returns a map of label to question-answer pair
func LoadDatasetLocal(datasetFilePath string) (map[string][]string, error) {
	file, err := os.Open(datasetFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dataset file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'       // Set the delimiter to pipe
	reader.LazyQuotes = true // Allow lazy quotes to handle bare quotes

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read dataset file: %v", err)
	}

	labelToQA := make(map[string][]string)
	for i, record := range records {
		if len(record) != 2 {
			log.Printf("Skipping invalid record at line %d: %v\n", i+1, record)
			continue
		}
		label := "LABEL_" + strconv.Itoa(i+1) // Adjust label to match dataset row numbers
		labelToQA[label] = record
	}
	return labelToQA, nil
}

func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
	var chat models.AIRequest

	err := json.NewDecoder(req.Body).Decode(&chat)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if chat.Prompt == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "masukin pertanyaan dulu ya kakak 🤗")
		return
	}

	client := resty.New()

	// Hugging Face API URL dan token
	apiUrl := modules.GetEnv("HUGGINGFACE_API_URL")
	apiToken := "Bearer " + tokenmodel

	response, err := client.R().
		SetHeader("Authorization", apiToken).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"inputs": "` + chat.Prompt + `"}`).
		Post(apiUrl)

	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	// Periksa jika model sedang dimuat
	if response.StatusCode() == http.StatusServiceUnavailable {
		utils.ErrorResponse(respw, req, http.StatusServiceUnavailable, "Internal Server Error", "Model sedang dimuat, coba lagi sebentar ya kakak 🙏 | HF Response: "+response.String())
		return
	}

	// Periksa jika model tidak ditemukan
	if response.StatusCode() == http.StatusNotFound {
		utils.ErrorResponse(respw, req, http.StatusNotFound, "Not Found", "Model tidak ditemukan | HF Response: "+response.String())
		return
	}

	// Periksa jika model mengembalikan status code lain
	if response.StatusCode() != http.StatusOK {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: "+response.String())
		return
	}

	// Handle the expected nested array structure
	var nestedData [][]map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &nestedData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error decoding response: "+err.Error()+" | Server HF Response: "+response.String())
		return
	}

	// Flatten the nested array structure
	var flatData []map[string]interface{}
	for _, d := range nestedData {
		flatData = append(flatData, d...)
	}

	// Log the generated text to the terminal
	log.Printf("Generated text: %v", flatData)

	// Extracting the highest scoring label from the model output
	var bestLabel string
	var highestScore float64
	for _, item := range flatData {
		label, labelOk := item["label"].(string)
		score, scoreOk := item["score"].(float64)
		if labelOk && scoreOk && (bestLabel == "" || score > highestScore) {
			bestLabel = label
			highestScore = score
		}
	}

	if bestLabel != "" {
		// Load the dataset from local file
		datasetFilePath := "../rf1.csv"
		labelToQA, err := LoadDatasetLocal(datasetFilePath)
		if err != nil {
			utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "server error: could not load dataset: "+err.Error())
			return
		}

		// Get the answer corresponding to the best label
		record, ok := labelToQA[bestLabel]
		if !ok {
			utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "server error: label not found in dataset")
			return
		}

		answer := record[1]

		utils.WriteJSON(respw, http.StatusOK, map[string]string{
			"prompt":   chat.Prompt,
			"response": answer,
			"label":    bestLabel,
			"score":    strconv.FormatFloat(highestScore, 'f', -1, 64),
		})
	} else {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : response")
	}
}
