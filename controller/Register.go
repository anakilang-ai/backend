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
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// Logger instance
var log = logrus.New()

func init() {
	// Customize the logger if needed
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

// SignUp handles user registration requests
func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	// Decode the request body into the user struct
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		log.WithError(err).Error("Failed to parse request body")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
		return
	}

	// Validate input fields
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		log.Warn("Incomplete data in signup request")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Please provide all required fields")
		return
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		log.WithField("email", user.Email).Warn("Invalid email format")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Invalid email format")
		return
	}

	// Check if the email is already registered
	userExists, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		log.WithError(err).Error("Error checking if user exists")
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Error checking if user exists: "+err.Error())
		return
	}
	if userExists.Email != "" {
		log.WithField("email", user.Email).Warn("Email already registered")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email is already registered")
		return
	}

	// Validate password
	if strings.Contains(user.Password, " ") {
		log.Warn("Password contains spaces")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password cannot contain spaces")
		return
	}
	if len(user.Password) < 8 {
		log.Warn("Password too short")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password must be at least 8 characters long")
		return
	}
	if user.Password != user.Confirmpassword {
		log.Warn("Password and confirmation password do not match")
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Passwords do not match")
		return
	}

	// Generate salt and hash the password
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.WithError(err).Error("Error generating salt")
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Error generating salt: "+err.Error())
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

	// Insert the new user into the database
	insertedID, err := utils.InsertOneDoc(db, col, userData)
	if err != nil {
		log.WithError(err).Error("Error inserting new user")
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Error inserting user: "+err.Error())
		return
	}

	// Respond with success
	resp := map[string]interface{}{
		"message":    "Registration successful",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	log.WithField("email", user.Email).Info("User registered successfully")
	utils.WriteJSON(respw, http.StatusCreated, resp)
}
