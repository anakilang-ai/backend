package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/anakilang-ai/backend/utils"
	model "github.com/anakilang-ai/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}
	userExists, _ := utils.GetUserFromEmail(user.Email, db)
	if user.Email == userExists.Email {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email sudah terdaftar")
		return
	}
	if strings.Contains(user.Password, " ") {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi")
		return
	}
	if len(user.Password) < 8 {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter")
		return
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	userData := bson.M{
		"namalengkap": user.NamaLengkap,
		"email":       user.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"salt":        hex.EncodeToString(salt),
	}
	insertedID, err := utils.InsertOneDoc(db, col, userData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : insert data, "+err.Error())
		return
	}
	resp := map[string]any{
		"message":    "berhasil mendaftar",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	utils.WriteJSON(respw, http.StatusCreated, resp)
}