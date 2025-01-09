package certs

import (
	. "github.com/Obsec24/control-master/common"
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
		Log("Unpacking the uploaded apk")
		Command(nil, "apktool", "d", "-s", "base.apk", "-o", "apk")
		send, _ := json.Marshal(resp{true, "File uploaded successfuly"})
		fmt.Fprintf(res, string(send))
	} else {
		send, _ := json.Marshal(resp{false, "Error uploading file"})
		fmt.Fprintf(res, string(send))
	}
}

func info(res http.ResponseWriter, req *http.Request) {
	if _, err := os.Stat("base.apk"); os.IsNotExist(err) {
		send, _ := json.Marshal(resp{false, "Upload APK to be analyzed"})
		fmt.Fprintf(res, string(send))
	} else {
		Log("Getting certificate information")
		Command(res, "./cert_info.py", "apk")
	}
}
