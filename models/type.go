package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the application.
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FullName        string             `bson:"full_name,omitempty" json:"full_name,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	ConfirmPassword string             `bson:"confirm_password,omitempty" json:"confirm_password,omitempty"`
	Salt            string             `bson:"salt,omitempty" json:"salt,omitempty"`
}

// Password represents a password change request.
type Password struct {
	Password        string `bson:"password,omitempty" json:"password,omitempty"`
	NewPassword     string `bson:"new_password,omitempty" json:"new_password,omitempty"`
	ConfirmPassword string `bson:"confirm_password,omitempty" json:"confirm_password,omitempty"`
}

// AIRequest represents a request to the AI.
type AIRequest struct {
	Prompt    string    `bson:"prompt,omitempty" json:"prompt,omitempty"`
	AIResp    string    `bson:"ai_resp,omitempty" json:"ai_resp,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// AIResponse represents a response from the AI.
type AIResponse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AIRequest AIRequest          `bson:"ai_request,omitempty" json:"ai_request,omitempty"`
	Response  string             `bson:"response,omitempty" json:"response,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// Credential represents authentication credentials.
type Credential struct {
	Status  int    `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// Response represents a generic API response.
type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// Payload represents the payload of a token.
type Payload struct {
	ID    primitive.ObjectID `json:"id" bson:"id"`
	Email string             `json:"email" bson:"email"`
	Exp   time.Time          `json:"exp" bson:"exp"`
	Iat   time.Time          `json:"iat" bson:"iat"`
	Nbf   time.Time          `json:"nbf" bson:"nbf"`
}
