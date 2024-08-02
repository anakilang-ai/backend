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

// DBInfo holds the database connection information.
type DBInfo struct {
	DBString string // Connection string for MongoDB
	DBName   string // Name of the MongoDB database
}

// MongoConnect establishes a connection to MongoDB and returns a database instance.
func MongoConnect(mconn DBInfo) (*mongo.Database, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(mconn.DBName), nil
}

// InsertOneDoc inserts a single document into the specified collection and returns the inserted ID.
func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (primitive.ObjectID, error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("inserted ID is not of type primitive.ObjectID")
	}

	return id, nil
}

// GetUserFromEmail retrieves a user document by email.
func GetUserFromEmail(email string, db *mongo.Database) (models.User, error) {
	collection := db.Collection("users")
	filter := bson.M{"email": email}
	var doc models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email not found: %w", err)
		}
		return doc, fmt.Errorf("server error while retrieving user: %w", err)
	}

	return doc, nil
}

// GetAllDocs retrieves all documents from a specified collection that match the given filter.
func GetAllDocs[T any](db *mongo.Database, col string, filter bson.M) ([]T, error) {
	ctx := context.TODO()
	collection := db.Collection(col)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []T
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return docs, nil
}

// GetUserFromID retrieves a user document by its ID.
func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (models.User, error) {
	collection := db.Collection("users")
	filter := bson.M{"_id": _id}
	var doc models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s: %w", _id.Hex(), err)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %w", _id.Hex(), err)
	}

	return doc, nil
}
