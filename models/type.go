package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name            string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	ConfirmPassword string             `bson:"confirmpassword,omitempty" json:"confirmpassword,omitempty"`
	Salt            string             `bson:"salt,omitempty" json:"salt,omitempty"`
}

// Password represents a password change request
type Password struct {
	CurrentPassword string `bson:"password,omitempty" json:"password,omitempty"`
	NewPassword     string `bson:"newpass,omitempty" json:"newpassword,omitempty"`
	ConfirmPassword string `bson:"confirmpass,omitempty" json:"confirmpassword,omitempty"`
}

// AIRequest represents a request to an AI service
type AIRequest struct {
	Prompt     string    `bson:"prompt,omitempty" json:"prompt,omitempty"`
	AIResponse string    `bson:"airesp,omitempty" json:"airesp,omitempty"`
	CreatedAt  time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// AIResponse represents the response from an AI service
type AIResponse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AIRequest AIRequest          `bson:"airequest,omitempty" json:"airequest,omitempty"`
	Response  string             `bson:"response,omitempty" json:"response,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// Credential represents an authentication token and its status
type Credential struct {
	Status  int    `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// Response represents a generic API response
type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

// Payload represents a JWT token payload
type Payload struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}
