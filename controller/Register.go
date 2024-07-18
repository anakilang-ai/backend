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

func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
		return
	}

	// Validate input
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Please complete all required fields")
		return
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Invalid email format")
		return
	}

	userExists, _ := utils.GetUserFromEmail(user.Email, db)
	if user.Email == userExists.Email {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email is already registered")
		return
	}

	if strings.Contains(user.Password, " ") {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password should not contain spaces")
		return
	}

	if len(user.Password) < 8 {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password should be at least 8 characters long")
		return
	}

	if user.Password != user.Confirmpassword {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Passwords do not match")
		return
	}

	// Generate salt
	salt := make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error: unable to generate salt")
		return
	}

	// Hash password
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// Prepare user data
	userData := bson.M{
		"namalengkap": user.NamaLengkap,
		"email":       user.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"salt":        hex.EncodeToString(salt),
	}

	// Insert user data into database
	insertedID, err := utils.InsertOneDoc(db, col, userData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Server error: unable to insert data, "+err.Error())
		return
	}

	// Prepare response
	resp := map[string]any{
		"message":    "Successfully registered",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}

	// Send response
	utils.WriteJSON(respw, http.StatusCreated, resp)
}
