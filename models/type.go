package models

//mengimpor dua package
import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// merepresentasikan data pengguna yang akan disimpan di database MongoDB dan diekspos melalui JSON.
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap     string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	Confirmpassword string             `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
	Salt            string             `bson:"salt,omitempty" json:"salt,omitempty"`
}

// merepresentasikan data yang terkait dengan password pengguna,  bisa untuk update password atau saat registrasi, yang akan disimpan di database MongoDB dan diekspos melalui JSON.
type Password struct {
	Password        string `bson:"password,omitempty" json:"password,omitempty"`
	Newpassword     string `bson:"newpass,omitempty" json:"newpass,omitempty"`
	Confirmpassword string `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
}

// merepresentasikan data permintaan yang diajukan ke AI dalam program ini.
type AIRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User      User               `bson:"user,omitempty" json:"user,omitempty"`
	Query     string             `bson:"query,omitempty" json:"query,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type AIResponse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AIRequest AIRequest          `bson:"airequest,omitempty" json:"airequest,omitempty"`
	Response  string             `bson:"response,omitempty" json:"response,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type Credential struct {
	Status  int    `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}
