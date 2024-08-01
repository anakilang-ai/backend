package controller

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/anakilang-ai/backend/utils"
	model "github.com/anakilang-ai/backend/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// LogIn handles user login, validates credentials, and generates an authentication token.
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privateKey string) {
	var user model.User

	// Decode the request body into the user struct
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body: "+err.Error())
		return
	}

	// Validate the user input
	if user.Email == "" || user.Password == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// Validate email format
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// Fetch user from the database
	existingUser, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: "+err.Error())
		return
	}

	// Decode salt and hash the password
	salt, err := hex.DecodeString(existingUser.Salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: decoding salt")
		return
	}

	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existingUser.Password {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "password salah")
		return
	}

	// Generate authentication token
	tokenString, err := utils.Encode(existingUser.ID, existingUser.Email, privateKey)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: generating token")
		return
	}

	// Send successful response with token
	resp := map[string]string{
		"status":  "success",
		"message": "login berhasil",
		"token":   tokenString,
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}
