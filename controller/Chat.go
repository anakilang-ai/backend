package controller

//mengimpor beberapa package yang dibutuhkan untuk program Go, termasuk package dari repository GitHub eksternal.
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
	// Fungsi Chat menangani request HTTP untuk chat dengan AI

	// Definisikan variable chat untuk menyimpan data request
	var chat models.AIRequest

	// Decode request body ke dalam struct chat
	err := json.NewDecoder(req.Body).Decode(&chat)
	if err != nil {
		// Jika terjadi error saat parsing request body, kembalikan Bad Request error
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	// Validasi query, jika kosong kembalikan Bad Request error
	if chat.Query == "" {
		helper.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// Inisialisasi client untuk request ke Hugging Face API
	client := resty.New()

	// Ambil Hugging Face API URL dan token dari environment variable
	apiUrl := modules.GetEnv("HUGGINGFACE_API_KEY")
	apiToken := "Bearer " + tokenmodel

	// Variable untuk menyimpan response dan error dari API
	var response *resty.Response
	var retryCount int
	maxRetries := 5
	retryDelay := 20 * time.Second

	// Looping untuk melakukan request ke Hugging Face API dengan mekanisme retry
	for retryCount < maxRetries {
		// Set header untuk Authorization dan Content-Type
		response, err = client.R().
			SetHeader("Authorization", apiToken).
			SetHeader("Content-Type", "application/json").
			SetBody(`{"inputs": "` + chat.Query + `"}`).
			Post(apiUrl)

		if err != nil {
			// Jika terjadi error saat request, log error dan lanjutkan ke retry
			log.Fatalf("Error making request: %v", err)
		}

		// Cek status code response
		if response.StatusCode() == http.StatusOK {
			// Jika status code 200 (OK), break dari loop retry
			break
		} else {
			// Jika status code bukan 200, coba parsing error response
			var errorResponse map[string]interface{}
			err = json.Unmarshal(response.Body(), &errorResponse)
			if err == nil && errorResponse["error"] == "Model is currently loading" {
				// Jika error menunjukkan model sedang loading, lakukan retry
				retryCount++
				time.Sleep(retryDelay)
				continue
			}
			// Jika error lain atau gagal parsing error response, kembalikan Internal Server Error
			helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error from Hugging Face API "+string(response.Body()))
			return
		}
	}

	// Cek status code response setelah loop retry
	if response.StatusCode() != 200 {
		// Jika status code tetap bukan 200 setelah retry, kembalikan Internal Server Error
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error from Hugging Face API "+string(response.Body()))
		return
	}

	// Decode response body ke dalam slice map
	var data []map[string]interface{}

	err = json.Unmarshal(response.Body(), &data)
	if err != nil {
		// Jika terjadi error saat parsing response body, kembalikan Internal Server Error
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error parsing response body "+err.Error())
		return
	}

	// Proses data response
	if len(data) > 0 {
		// Ekstrak generated text dari response
		generatedText, ok := data[0]["generated_text"].(string)
		if !ok {
			// Jika gagal ekstrak generated text, kembalikan Internal Server Error
			helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "error extracting generated text")
			return
		}
		// Kirim response sukses dengan generated text
		helper.WriteJSON(respw, http.StatusOK, map[string]string{"answer": generatedText})
	} else {
		// Jika tidak ada data response, kembalikan Internal Server Error
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: response")
	}
}
