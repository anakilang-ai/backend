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

// LogIn handles user login, validates credentials, and returns a JWT token on success.
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User

	// Decode request body
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
		return
	}

	// Check for empty email or password
	if user.Email == "" || user.Password == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Please provide both email and password")
		return
	}

	// Validate email format
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Invalid email format")
		return
	}

	// Retrieve user document by email
	existsDoc, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error: could not retrieve user by email: "+err.Error())
		return
	}

	// Decode user's stored salt
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error: could not decode salt")
		return
	}

	// Hash the provided password with the stored salt
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "Incorrect password")
		return
	}

	// Generate JWT token
	tokenString, err := utils.Encode(user.ID, user.Email, privatekey)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error: could not generate token")
		return
	}

	// Successful login response
	resp := map[string]string{
		"status":  "success",
		"message": "Login successful",
		"token":   tokenString,
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}
