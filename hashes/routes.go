package hashes

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/hashes", hashes)
}
