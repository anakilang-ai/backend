package modules

// digunakan untuk mengimpor package utils dari repository GitHub anakilang-ai/backend.
import (
	"github.com/anakilang-ai/backend/utils"
)

// mendefinisikan variabel MongoString dengan tipe data string dan menginisialisasi nilainya dengan hasil dari fungsi GetEnv yang memiliki parameter "MONGOSTRING".
var MongoString string = GetEnv("MONGOSTRING")

var mongoinfo = utils.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

var Mongoconn, ErrorMongoconn = utils.MongoConnect(mongoinfo)
