package libraries

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/literadar", libScout)
	http.HandleFunc("/libscout", liteRadar)
}
