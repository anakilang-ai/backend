package helper

//mengimpor beberapa package yang dibutuhkan untuk program Go.
import (
	"encoding/json"
	"log"
	"net/http"
)

// mengirimkan response error ke client (browser atau aplikasi lain yang mengirimkan request)  dalam format JSON.
func ErrorResponse(respw http.ResponseWriter, req *http.Request, statusCode int, err, msg string) {
	resp := map[string]string{
		"error":   err,
		"message": msg,
	}
	WriteJSON(respw, statusCode, resp)
}

// menuliskan data apapun (represented by content parameter) ke response HTTP dalam format JSON.
func WriteJSON(respw http.ResponseWriter, statusCode int, content any) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

// mengencode data apapun (represented by strc parameter) ke format JSON string.
func Jsonstr(strc any) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}
