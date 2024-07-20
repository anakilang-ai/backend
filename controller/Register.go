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

	// Decode request body
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
		return
	}

	// Validate user input
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Mohon untuk melengkapi data")
		return
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email tidak valid!")
		return
	}

	userExists, _ := utils.GetUserFromEmail(user.Email, db)
	if userExists.Email != "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email sudah terdaftar")
		return
	}

	if strings.Contains(user.Password, " ") {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password tidak boleh mengandung spasi")
		return
	}

	if len(user.Password) < 8 {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password minimal 8 karakter")
		return
	}

	if user.Password != user.Confirmpassword {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Password dan konfirmasi password tidak cocok")
		return
	}

	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: salt")
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

	// Insert user data into the database
	insertedID, err := utils.InsertOneDoc(db, col, userData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: insert data, "+err.Error())
		return
	}

	// Send response
	resp := map[string]any{
		"message":    "Berhasil mendaftar",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	utils.WriteJSON(respw, http.StatusCreated, resp)
}
