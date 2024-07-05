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

// user
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body " + err.Error())
		return
	}
	if user.Email == "" || user.Password == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}
	existsDoc, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : get email " + err.Error())
		return
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "password salah")
		return
	}
	tokenstring, err := utils.Encode(user.ID, user.Email, privatekey)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : token")
		return
	}
	resp := map[string]string{
		"status":  "success",
		"message": "login berhasil",
		"token":   tokenstring,
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}