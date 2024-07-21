package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}

func Encode(id primitive.ObjectID, email, privateKey string) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.Set("id", id)
	token.SetString("email", email)
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create secret key: %v", err)
	}
	return token.V4Sign(secretKey, nil), nil
}

func Decode(publicKey string, tokenstring string) (payload Payload, err error) {
	var token *paseto.Token
	pubKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey)
	if err != nil {
		return payload, fmt.Errorf("failed to create public key: %v", err)
	}
	parser := paseto.NewParser()
	token, err = parser.ParseV4Public(pubKey, tokenstring, nil)
	if err != nil {
		return payload, fmt.Errorf("failed to parse token: %v", err)
	}
	err = json.Unmarshal(token.ClaimsJSON(), &payload)
	if err != nil {
		return payload, fmt.Errorf("failed to unmarshal token claims: %v", err)
	}
	return payload, nil
}

func GenerateKey() (privateKey, publicKey string) {
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey = secretKey.Public().ExportHex()
	privateKey = secretKey.ExportHex()
	return privateKey, publicKey
}
