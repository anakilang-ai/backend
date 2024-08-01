package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payload represents the structure of the PASETO token payload.
type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}

// Encode creates a PASETO token with the given payload and private key.
func Encode(id primitive.ObjectID, email, privateKey string) (string, error) {
	token := paseto.NewToken()
	now := time.Now()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(2 * time.Hour))
	token.Set("id", id.Hex())
	token.SetString("email", email)
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create secret key from hex: %w", err)
	}
	return token.V4Sign(secretKey, nil), nil
}

// Decode verifies and parses a PASETO token string using the provided public key.
func Decode(publicKey string, tokenString string) (Payload, error) {
	var payload Payload
	pubKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey)
	if err != nil {
		return payload, fmt.Errorf("failed to create public key from hex: %w", err)
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(pubKey, tokenString, nil)
	if err != nil {
		return payload, fmt.Errorf("failed to parse token: %w", err)
	}

	err = json.Unmarshal(token.ClaimsJSON(), &payload)
	if err != nil {
		return payload, fmt.Errorf("failed to unmarshal token claims: %w", err)
	}

	return payload, nil
}

// GenerateKey generates a new pair of PASETO private and public keys.
func GenerateKey() (privateKey, publicKey string) {
	secretKey := paseto.NewV4AsymmetricSecretKey() // Don't share this!!!
	publicKey = secretKey.Public().ExportHex()     // DO share this one
	privateKey = secretKey.ExportHex()
	return privateKey, publicKey
}
