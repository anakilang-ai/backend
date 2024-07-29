package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	model "github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/utils"
	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// SignUp handles user registration
func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	// Decode the request body into the user model
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body: "+err.Error())
		return
	}

	// Validate that all required fields are present
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	// Validate email format
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}

	// Check if the email is already registered
	existingUser, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: error checking user existence, "+err.Error())
		return
	}
	if existingUser.Email != "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email sudah terdaftar")
		return
	}

	// Validate the password
	if strings.Contains(user.Password, " ") {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi")
		return
	}
	if len(user.Password) < 8 {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter")
		return
	}
	if user.Password != user.Confirmpassword {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password dan konfirmasi password tidak cocok")
		return
	}

	// Generate salt and hash the password
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: error generating salt, "+err.Error())
		return
	}
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// Prepare the user data for insertion
	userData := bson.M{
		"namalengkap": user.NamaLengkap,
		"email":       user.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"salt":        hex.EncodeToString(salt),
	}

	// Insert user data into the database
	insertedID, err := utils.InsertOneDoc(db, col, userData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server: error inserting data, "+err.Error())
		return
	}

	// Respond with success
	resp := map[string]interface{}{
		"message":    "berhasil mendaftar",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	utils.WriteJSON(respw, http.StatusCreated, resp)
}
