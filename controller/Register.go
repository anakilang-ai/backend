package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"context"
	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/argon2"
	"time"
)

var dbClient *mongo.Client
var dbName = "yourdbname"
var colName = "users"

// User model
type User struct {
	ID              string `json:"id" bson:"_id,omitempty"`
	NamaLengkap     string `json:"namalengkap"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Confirmpassword string `json:"confirmpassword,omitempty"`
	Salt            string `json:"salt"`
}

// ConnectDB establishes a connection to the MongoDB database
func ConnectDB() {
	var err error
	dbClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

// GetUserFromEmail fetches a user by email from the database
func GetUserFromEmail(email string, db *mongo.Database) (User, error) {
	var user User
	collection := db.Collection(colName)
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return user, err
}

// ErrorResponse sends a JSON error response
func ErrorResponse(respw http.ResponseWriter, req *http.Request, status int, title string, detail string) {
	resp := map[string]string{
		"status":  title,
		"message": detail,
	}
	WriteJSON(respw, status, resp)
}

// WriteJSON sends a JSON response
func WriteJSON(respw http.ResponseWriter, status int, data interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(status)
	json.NewEncoder(respw).Encode(data)
}

// InsertOneDoc inserts a document into the specified collection
func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (string, error) {
	collection := db.Collection(col)
	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}
	id := result.InsertedID.(bson.ObjectId).Hex()
	return id, nil
}

// SignUp handles user registration
func SignUp(db *mongo.Database, respw http.ResponseWriter, req *http.Request) {
	var user User

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}

	if user.NamaLengkap == "" || user.Email == "" || user.Password == "" || user.Confirmpassword == "" {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid!")
		return
	}
	userExists, _ := GetUserFromEmail(user.Email, db)
	if user.Email == userExists.Email {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email sudah terdaftar")
		return
	}
	if strings.Contains(user.Password, " ") {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password tidak boleh mengandung spasi")
		return
	}
	if len(user.Password) < 8 {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "password minimal 8 karakter")
		return
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	userData := bson.M{
		"namalengkap": user.NamaLengkap,
		"email":       user.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"salt":        hex.EncodeToString(salt),
	}
	insertedID, err := InsertOneDoc(db, colName, userData)
	if err != nil {
		ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : insert data, "+err.Error())
		return
	}
	resp := map[string]any{
		"message":    "berhasil mendaftar",
		"insertedID": insertedID,
		"data": map[string]string{
			"email": user.Email,
		},
	}
	WriteJSON(respw, http.StatusCreated, resp)
}

func main() {
	ConnectDB()

	r := mux.NewRouter()
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		SignUp(dbClient.Database(dbName), w, r)
	}).Methods("POST")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
