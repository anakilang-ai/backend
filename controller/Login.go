package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"context"
	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/argon2"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var dbClient *mongo.Client
var dbName = "yourdbname"
var colName = "users"
var privateKey = "yourprivatekey"

// User model
type User struct {
	ID              string `json:"id" bson:"_id,omitempty"`
	NamaLengkap     string `json:"namalengkap"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Confirmpassword string `json:"confirmpassword,omitempty"`
	Salt            string `json:"salt"`
}

// JWT Claims struct
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
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

// Encode generates a JWT token
func Encode(userID, email, privatekey string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(privatekey))
	return tokenString, err
}

// LogIn handles user login
func LogIn(db *mongo.Database, respw http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "error parsing request body "+err.Error())
		return
	}
	if user.Email == "" || user.Password == "" {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "mohon untuk melengkapi data")
		return
	}
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		ErrorResponse(respw, req, http.StatusBadRequest, "Bad Request", "email tidak valid")
		return
	}
	existsDoc, err := GetUserFromEmail(user.Email, db)
	if err != nil {
		ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : get email "+err.Error())
		return
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : salt")
		return
	}
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		ErrorResponse(respw, req, http.StatusUnauthorized, "Unauthorized", "password salah")
		return
	}
	tokenstring, err := Encode(user.ID, user.Email, privateKey)
	if err != nil {
		ErrorResponse(respw, req, http.StatusInternalServerError, "Internal Server Error", "kesalahan server : token")
		return
	}
	resp := map[string]string{
		"status":  "success",
		"message": "login berhasil",
		"token":   tokenstring,
	}
	WriteJSON(respw, http.StatusOK, resp)
}

func main() {
	ConnectDB()

	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LogIn(dbClient.Database(dbName), w, r)
	}).Methods("POST")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
