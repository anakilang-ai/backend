package modules

// digunakan untuk mengimpor package utils dari repository GitHub anakilang-ai/backend.
import (
	"github.com/anakilang-ai/backend/utils"
)

// mendefinisikan variabel MongoString dengan tipe data string dan menginisialisasi nilainya dengan hasil dari fungsi GetEnv yang memiliki parameter "MONGOSTRING".
var MongoString string = GetEnv("MONGOSTRING")

// mendefinisikan variabel mongoinfo dengan tipe data utils.DBInfo dan menginisialisasi nilainya dengan struktur utils.DBInfo
var mongoinfo = utils.DBInfo{
	DBString: MongoString,
	DBName:   "ailang",
}

// membuat koneksi ke database MongoDB dan menyimpan hasil koneksi dalam variabel Mongoconn.
var Mongoconn, ErrorMongoconn = utils.MongoConnect(mongoinfo)
