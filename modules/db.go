package modules

import (
	"log"

	"github.com/anakilang-ai/backend/utils"
)

// Initialize MongoDB connection string from environment variable.
var MongoString = GetEnv("MONGOSTRING")

// Define the MongoDB connection information.
var mongoinfo = utils.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

// Initialize MongoDB connection and handle any connection error.
var (
	Mongoconn, ErrorMongoconn = utils.MongoConnect(mongoinfo)
)

func init() {
	if ErrorMongoconn != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", ErrorMongoconn)
	} else {
		log.Println("Successfully connected to MongoDB")
	}
}
