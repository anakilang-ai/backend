// db.go
package modules

import (
	"log"

	"github.com/anakilang-ai/backend/config" // Import config package for environment variables
	"github.com/anakilang-ai/backend/utils"
)

// Initialize MongoDB connection details
var (
	MongoString = config.GetEnv("MONGOSTRING")
	mongoInfo   = utils.DBInfo{
		DBString: MongoString,
		DBName:   "ailang",
	}
	MongoConn, MongoConnErr = utils.MongoConnect(mongoInfo)
)

func init() {
	if MongoConnErr != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", MongoConnErr)
	}
	log.Println("Successfully connected to MongoDB")
}
