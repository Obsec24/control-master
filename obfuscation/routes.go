package obfuscation

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/apkid", apkid)
}
