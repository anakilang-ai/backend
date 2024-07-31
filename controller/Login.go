package controller

import (
	"encoding/json"
	"net/http"

	model "github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/utils"
	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

func Login(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var loginData model.Login

	err := json.NewDecoder(req.Body).Decode(&loginData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}
	if err := checkmail.ValidateFormat(loginData.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	userExists, err := utils.GetUserFromEmail(loginData.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "email atau password salah")
		return
	}

	salt, err := hex.DecodeString(userExists.Salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : decode salt")
		return
	}
	hashedPassword := argon2.IDKey([]byte(loginData.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hashedPassword) != userExists.Password {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "email atau password salah")
		return
	}

	resp := map[string]any{
		"message": "berhasil login",
		"data": map[string]string{
			"email": userExists.Email,
		},
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}
