package modules

import (
	"log"

	"github.com/anakilang-ai/backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoString holds the MongoDB connection string from environment variables.
var MongoString = utils.GetEnv("MONGOSTRING")

// Define the MongoDB connection information.
var mongoinfo = utils.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

// MongoDB connection and error variables.
var (
	Mongoconn      *mongo.Database
	ErrorMongoconn error
)

func init() {
	var err error
	Mongoconn, err = utils.MongoConnect(mongoinfo)
	if err != nil {
		ErrorMongoconn = err
		log.Fatalf("Failed to connect to MongoDB: %v", ErrorMongoconn)
	} else {
		log.Println("Successfully connected to MongoDB")
	}
}
