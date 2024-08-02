package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payload represents the JWT payload structure.
type Payload struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}

// Encode generates a PASETO token with the given ID, email, and private key.
func Encode(id primitive.ObjectID, email, privateKey string) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.Set("id", id)
	token.SetString("email", email)

	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create secret key from hex: %v", err)
	}

	return token.V4Sign(secretKey, nil), nil
}

// Decode parses a PASETO token and returns the payload.
func Decode(publicKey string, tokenString string) (Payload, error) {
	var payload Payload
	var pubKey paseto.V4AsymmetricPublicKey

	pubKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey)
	if err != nil {
		return payload, fmt.Errorf("failed to create public key from hex: %v", err)
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(pubKey, tokenString, nil)
	if err != nil {
		return payload, fmt.Errorf("failed to parse token: %v", err)
	}

	err = json.Unmarshal(token.ClaimsJSON(), &payload)
	if err != nil {
		return payload, fmt.Errorf("failed to unmarshal token claims: %v", err)
	}

	return payload, nil
}

// GenerateKey generates a new pair of PASETO asymmetric keys.
func GenerateKey() (privateKey, publicKey string) {
	secretKey := paseto.NewV4AsymmetricSecretKey() // DO NOT share this!!!
	publicKey = secretKey.Public().ExportHex()     // Safe to share
	privateKey = secretKey.ExportHex()
	return privateKey, publicKey
}
