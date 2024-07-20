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

// LogIn menangani proses login pengguna, termasuk validasi kredensial dan pembuatan token JWT
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request, privatekey string) {
	var user model.User

	// Mengurai body permintaan menjadi objek user
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Kesalahan dalam parsing body permintaan: "+err.Error())
		return
	}

	// Validasi email dan password
	if user.Email == "" || user.Password == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Mohon untuk melengkapi data")
		return
	}

	// Validasi format email
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email tidak valid")
		return
	}

	// Mendapatkan dokumen pengguna dari email
	existsDoc, err := utils.GetUserFromEmail(user.Email, db)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: gagal mendapatkan email: "+err.Error())
		return
	}

	// Dekode salt dari string heksadesimal
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: gagal mendekode salt")
		return
	}

	// Hash password menggunakan Argon2 dan salt yang diambil dari database
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		utils.ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "Password salah")
		return
	}

	// Membuat token JWT
	tokenstring, err := utils.Encode(existsDoc.ID, user.Email, privatekey)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: gagal membuat token")
		return
	}

	// Mengirim respon sukses dengan token JWT
	resp := map[string]string{
		"status":  "success",
		"message": "Login berhasil",
		"token":   tokenstring,
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}
