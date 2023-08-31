package model

import (
	"net/http"
)

func AddCORSHeader(w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json")
	//Necessary for request cords
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}
