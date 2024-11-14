package traffic

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/config", config)
	http.HandleFunc("/cert", cert)
	http.HandleFunc("/hook", hooker)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/phase-one", phaseOne)
	http.HandleFunc("/phase-two", phaseTwo)
	http.HandleFunc("/analysis", analysis)
	http.HandleFunc("/result", result)
	http.HandleFunc("/raw-phase-one", rawPhaseOne)
	http.HandleFunc("/raw-phase-two", rawPhaseTwo)
	http.HandleFunc("/screenshot-phase-one", screenshotPhaseOne)
	http.HandleFunc("/screenshot-phase-two", screenshotPhaseTwo)
	http.HandleFunc("/sanitize", sanitize)
}
