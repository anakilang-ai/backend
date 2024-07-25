package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/websocket"
)

// Message represents a chat message
type Message struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatRoom represents a chat room
type ChatRoom struct {
	ID       string    `json:"id" bson:"_id,omitempty"`
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
}

var (
	dbClient *mongo.Client
	dbName   = "chatdb"
	colName  = "chatrooms"
)

// ConnectDB establishes a connection to the MongoDB database
func ConnectDB() {
	var err error
	dbClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = dbClient.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = dbClient.Ping(nil, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

// CreateChatRoom handles the creation of a new chat room
func CreateChatRoom(w http.ResponseWriter, r *http.Request) {
	var chatRoom ChatRoom
	err := json.NewDecoder(r.Body).Decode(&chatRoom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := dbClient.Database(dbName).Collection(colName)
	_, err = collection.InsertOne(r.Context(), chatRoom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chatRoom)
}

// GetChatRooms handles fetching all chat rooms
func GetChatRooms(w http.ResponseWriter, r *http.Request) {
	collection := dbClient.Database(dbName).Collection(colName)
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var chatRooms []ChatRoom
	err = cursor.All(r.Context(), &chatRooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatRooms)
}

// SendMessage handles sending a message to a chat room
func SendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message.Timestamp = time.Now()

	collection := dbClient.Database(dbName).Collection(colName)
	_, err = collection.UpdateOne(
		r.Context(),
		bson.M{"_id": roomID},
		bson.M{"$push": bson.M{"messages": message}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// WebSocketHandler handles websocket connections for real-time messaging
func WebSocketHandler(ws *websocket.Conn) {
	var message Message
	for {
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			log.Println("WebSocket error:", err)
			break
		}

		// Broadcast the message to all connected clients
		websocket.JSON.Send(ws, message)
	}
}

func main() {
	ConnectDB()

	r := mux.NewRouter()
	r.HandleFunc("/chatrooms", CreateChatRoom).Methods("POST")
	r.HandleFunc("/chatrooms", GetChatRooms).Methods("GET")
	r.HandleFunc("/chatrooms/{roomID}/messages", SendMessage).Methods("POST")
	r.Handle("/ws", websocket.Handler(WebSocketHandler))

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
