package main

import (
	"log"
	"net/http"

	//https://www.golangprograms.com/how-to-use-function-from-another-file-golang.html
	//"strconv"
	Config "rankapi/ConfigHelper"
	DataModel "rankapi/Model"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func defaultMsg(w http.ResponseWriter, r *http.Request) {

}

// go mod init rankapi permite crear el modgo
func main() {
	port := Config.ReadValue("Port")
	crt := Config.ReadValue("SSL_CRT_FILE")
	key := Config.ReadValue("SSL_KEY_FILE")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", defaultMsg)
	router.HandleFunc("/api/TierModels", DataModel.GetAllTiers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/TierModels/{id}", DataModel.GetTier).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/TierModels/{id}", DataModel.DeleteTier).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/TierModels", DataModel.CreateTier).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/TierModels/{id}", DataModel.UpdateTier).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/TierModels/{id}", DataModel.UpdateTier).Methods("PATCH", "OPTIONS")

	router.HandleFunc("/api/ItemModels/{id}/tierId/{tierId}/{ranking}", DataModel.AssignTier).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/ItemModels", DataModel.GetAllItems).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/ItemModels/{id}", DataModel.GetItem).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/ItemModels/{id}", DataModel.DeleteItem).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/ItemModels/ResetType/{id}", DataModel.ResetType).Methods("GET", "OPTIONS")

	//cors.AllowAll()
	handler := cors.AllowAll().Handler(router) //Para que no moleste la seguridad.

	//log.Fatal(http.ListenAndServe(":"+port, handler))
	log.Fatal(http.ListenAndServeTLS(":"+port, crt, key, handler))
}
