package utils

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

// DBInfo contains the connection string and database name
type DBInfo struct {
	DBString string
	DBName   string
}

// MongoConnect establishes a connection to the MongoDB database
func MongoConnect(mconn DBInfo) (*mongo.Database, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	return client.Database(mconn.DBName), nil
}

// InsertOneDoc inserts a document into the specified collection and returns the inserted ID
func InsertOneDoc(db *mongo.Database, col string, doc any) (primitive.ObjectID, error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %w", err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetUserFromEmail retrieves a user document from the "users" collection by email
func GetUserFromEmail(email string, db *mongo.Database) (models.User, error) {
	var user models.User
	collection := db.Collection("users")
	filter := bson.M{"email": email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, fmt.Errorf("email not found")
		}
		return user, fmt.Errorf("server error: %w", err)
	}
	return user, nil
}

// GetAllDocs retrieves all documents matching the filter from the specified collection
func GetAllDocs[T any](db *mongo.Database, col string, filter bson.M) (T, error) {
	var docs T
	collection := db.Collection(col)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return docs, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(context.TODO())
	if err = cursor.All(context.TODO(), &docs); err != nil {
		return docs, fmt.Errorf("failed to decode documents: %w", err)
	}
	return docs, nil
}

// GetUserFromID retrieves a user document from the "users" collection by ID
func GetUserFromID(id primitive.ObjectID, db *mongo.Database) (models.User, error) {
	var user models.User
	collection := db.Collection("users")
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, fmt.Errorf("no data found for ID %s", id.Hex())
		}
		return user, fmt.Errorf("error retrieving data for ID %s: %w", id.Hex(), err)
	}
	return user, nil
}
S