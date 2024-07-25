// db.go
package modules

import (
	"log"

	"github.com/anakilang-ai/backend/utils"
)

var MongoString string = GetEnv("MONGOSTRING")

var mongoinfo = utils.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

var Mongoconn, ErrorMongoconn = utils.MongoConnect(mongoinfo)

func init() {
	if ErrorMongoconn != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", ErrorMongoconn)
	} else {
		log.Println("Successfully connected to MongoDB")
	}
}
