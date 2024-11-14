package certs

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/info", info)
}
