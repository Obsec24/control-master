package libraries

import (
	. "control/common"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type resp struct {
	Ok  bool
	Msg string
}

func upload(res http.ResponseWriter, req *http.Request) {
	Log("Trying to upload APK")
	if err := UploadAPK(req); err == nil {
		send, _ := json.Marshal(resp{true, "File uploaded successfuly"})
		fmt.Fprintf(res, string(send))
	} else {
		send, _ := json.Marshal(resp{false, "Error uploading file"})
		fmt.Fprintf(res, string(send))
	}
}

func libScout(res http.ResponseWriter, req *http.Request) {
	if _, err := os.Stat("base.apk"); os.IsNotExist(err) {
		send, _ := json.Marshal(resp{false, "Upload APK to be analyzed"})
		fmt.Fprintf(res, string(send))
	} else {
		Log("Starting LibScout")
		Command(res, "./libScout.py", "base.apk", "LibScout/sdk", "LibScout/categorias-librerias-LibScout.json")
	}
}

func liteRadar(res http.ResponseWriter, req *http.Request) {
	if _, err := os.Stat("base.apk"); os.IsNotExist(err) {
		send, _ := json.Marshal(resp{false, "Upload APK to be analyzed"})
		fmt.Fprintf(res, string(send))
	} else {
		Log("Starting LiteRadar")
		Command(res, "./liteRadar.py", "base.apk")
	}
}
