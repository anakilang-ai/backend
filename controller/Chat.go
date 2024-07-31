package controller

import (
	"encoding/json"
	"net/http"
	"time"

	model "github.com/anakilang-ai/backend/models"
	"github.com/anakilang-ai/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SendMessage menangani pengiriman pesan
func SendMessage(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	var message model.Message

	err := json.NewDecoder(req.Body).Decode(&message)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if message.Sender == "" || message.Receiver == "" || message.Content == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}

	message.Timestamp = time.Now().UTC()
	messageData := bson.M{
		"sender":    message.Sender,
		"receiver":  message.Receiver,
		"content":   message.Content,
		"timestamp": message.Timestamp,
	}

	insertedID, err := utils.InsertOneDoc(db, col, messageData)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : insert data, "+err.Error())
		return
	}

	resp := map[string]any{
		"message":    "berhasil mengirim pesan",
		"insertedID": insertedID,
		"data":       messageData,
	}
	utils.WriteJSON(respw, http.StatusCreated, resp)
}

// GetMessages menangani pengambilan pesan antara dua pengguna
func GetMessages(db *mongo.Database, col string, respw http.ResponseWriter, req *http.Request) {
	sender := req.URL.Query().Get("sender")
	receiver := req.URL.Query().Get("receiver")

	if sender == "" || receiver == "" {
		utils.ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk menyediakan sender dan receiver")
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	}

	options := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := db.Collection(col).Find(req.Context(), filter, options)
	if err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : fetch data, "+err.Error())
		return
	}
	defer cursor.Close(req.Context())

	var messages []model.Message
	if err = cursor.All(req.Context(), &messages); err != nil {
		utils.ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : decode data, "+err.Error())
		return
	}

	resp := map[string]any{
		"messages": messages,
	}
	utils.WriteJSON(respw, http.StatusOK, resp)
}
