package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/utils"
	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

func SignUp(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var user models.User

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "Error parsing request body: "+err.Error())
	
