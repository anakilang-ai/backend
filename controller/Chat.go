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

// SignUp menangani pendaftaran pengguna baru, termasuk validasi data dan penyimpanan di database
func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user model.User

	// Mengurai body permintaan menjadi objek user
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Kesalahan dalam parsing body permintaan: "+err.Error())
		return
	}

	// Validasi data pengguna
	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Mohon untuk melengkapi data")
		return
	}

	// Validasi format email
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email tidak valid")
		return
	}

	// Periksa apakah email sudah terdaftar
	userExists, _ := utils.GetUserFromEmail(user.Email, db)
	if user.Email == userExists.Email {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Email sudah terdaftar")
		return
	}

	// Validasi password
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

	// Membuat salt untuk password hashing
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: gagal membuat salt")
		return
	}

	// Hash password menggunakan Argon2
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// Mempersiapkan data pengguna untuk disimpan di database
	userData := bson.M{
		"namalengkap": user.NamaLengkap,
		"email":       user.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"salt":        hex.EncodeToString(salt),
	}

	// Menyimpan data pengguna di database
	insertedID, err := utils.InsertOneDoc(db, col, userData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "Kesalahan server: gagal menyimpan data pengguna, "+err.Error())
		return
	}

	// Mengirim respon sukses
	resp := map[string]any{
		"message":    "Berhasil mendaftar",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	utils.WriteJSON(respw, http.StatusCreated, resp)
}
