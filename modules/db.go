package modules

import (
	"log"

	"github.com/anakilang-ai/backend/utils"
)

// MongoString holds the MongoDB connection string from environment variables.
var MongoString = GetEnv("MONGOSTRING")

// Define the MongoDB connection information.
var mongoinfo = utils.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

// MongoDB connection and error variables.
var (
	Mongoconn, ErrorMongoconn = utils.MongoConnect(mongoinfo)
)

func init() {
	// Log the outcome of MongoDB connection initialization.
	if ErrorMongoconn != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", ErrorMongoconn)
	} else {
		log.Println("Successfully connected to MongoDB")
	}
}
