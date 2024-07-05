package modules

import (
	"github.com/anakilang-ai/backend/helper"
)

var MongoString string = GetEnv("MONGOSTRING")

var mongoinfo = helper.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

var Mongoconn, ErrorMongoconn = helper.MongoConnect(mongoinfo)