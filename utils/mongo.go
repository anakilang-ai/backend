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
func MongoConnect(mconn DBInfo) (db *mongo.Database, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	return client.Database(mconn.DBName), nil
}

// InsertOneDoc inserts a single document into the specified collection and returns the inserted ID.
func InsertOneDoc(db *mongo.Database, col string, doc any) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %v", err)
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("inserted ID is not of type primitive.ObjectID")
	}
	return id, nil
}

// GetUserFromEmail retrieves a user document by email.
func GetUserFromEmail(email string, db *mongo.Database) (doc models.User, err error) {
	collection := db.Collection("users")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email not found: %v", err)
		}
		return doc, fmt.Errorf("server error while retrieving user: %v", err)
	}
	return doc, nil
}

// GetAllDocs retrieves all documents from a specified collection that match the given filter.
func GetAllDocs[T any](db *mongo.Database, col string, filter bson.M) (docs T, err error) {
	ctx := context.TODO()
	collection := db.Collection(col)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return docs, fmt.Errorf("failed to execute find query: %v", err)
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &docs)
	if err != nil {
		return docs, fmt.Errorf("failed to decode documents: %v", err)
	}
	return docs, nil
}

// GetUserFromID retrieves a user document by its ID.
func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc models.User, err error) {
	collection := db.Collection("users")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s: %v", _id.Hex(), err)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %v", _id.Hex(), err)
	}
	return doc, nil
}
