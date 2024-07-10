package utils

// memuat paket-paket yang diperlukan untuk mengelola konteks, menangani kesalahan, melakukan operasi input/output, mengakses model dalam aplikasi, dan berinteraksi dengan database MongoDB dalam proyek Go.
import (
	"context"
	"errors"
	"fmt"

	"github.com/anakilang-ai/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// menyimpan informasi koneksi database, yaitu string koneksi (DBString) dan nama database (DBName).
type DBInfo struct {
	DBString string
	DBName   string
}

// Fungsi MongoConnect digunakan untuk menghubungkan ke database MongoDB dan mengembalikan objek database.
func MongoConnect(mconn DBInfo) (db *mongo.Database, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return nil, err
	}
	return client.Database(mconn.DBName), nil
}

// memasukkan satu dokumen baru ke dalam koleksi tertentu pada database MongoDB.
func InsertOneDoc(db *mongo.Database, col string, doc any) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// mencari dan mengembalikan informasi pengguna dari database MongoDB berdasarkan email, dan memberikan pesan kesalahan yang sesuai jika pengguna tidak ditemukan atau terjadi kesalahan server.
func GetUserFromEmail(email string, db *mongo.Database) (doc models.User, err error) {
	collection := db.Collection("users")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetAllDocs[T any](db *mongo.Database, col string, filter bson.M) (docs T, err error) {
	// Fungsi GetAllDocs digunakan untuk mengambil semua dokumen dari collection database yang memenuhi filter tertentu.

	// Kontext background untuk operasi database
	ctx := context.TODO()
	// Ambil collection dari database
	collection := db.Collection(col)
	// Find dokumen yang sesuai dengan filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return
	}
	// Tutup cursor setelah selesai digunakan (defer)
	defer cursor.Close(ctx)
	// Decode semua dokumen ke dalam slice dengan tipe data generik T
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return
	}
	// Kembalikan slice berisi dokumen dan error (jika ada)
	return
}

func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc models.User, err error) {
	// Fungsi GetUserFromID digunakan untuk mengambil data user berdasarkan ID (_id) dari database.
	collection := db.Collection("users")
	// Filter untuk mencari dokumen dengan _id tertentu
	filter := bson.M{"_id": _id}
	// Find one document yang sesuai dengan filter
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}
