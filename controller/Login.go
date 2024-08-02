package controller

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	model "github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/utils"
	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// LogIn handles user login requests
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
		return
	}

	// Validate input
	if user.Email == "" || user.Password == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Please provide both email and password")
		return
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Invalid email format")
		return
	}

	// Retrieve user from database
	existsDoc, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error while retrieving user: "+err.Error())
		return
	}

	// Validate password
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error while decoding salt: "+err.Error())
		return
	}
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "Incorrect password")
		return
	}

	// Generate token
	tokenstring, err := utils.Encode(user.ID, user.Email, privatekey)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Error generating token: "+err.Error())
		return
	}

	// Send success response
	resp := map[string]string{
		"status":  "success",
		"message": "Login successful",
		"token":   tokenstring,
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}
