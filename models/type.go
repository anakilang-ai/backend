package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User merepresentasikan data pengguna dalam aplikasi
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap     string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	Confirmpassword string             `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
	Salt            string             `bson:"salt,omitempty" json:"salt,omitempty"`
}

// Password merepresentasikan struktur untuk mengubah password
type Password struct {
	Password        string `bson:"password,omitempty" json:"password,omitempty"`
	Newpassword     string `bson:"newpassword,omitempty" json:"newpassword,omitempty"`
	Confirmpassword string `bson:"confirmpassword,omitempty" json:"confirmpassword,omitempty"`
}

// AIRequest me
