package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var router *mux.Router = mux.NewRouter()

	/** Router handlers */
	router.HandleFunc("/api/ping", pingAPI).Methods("GET")

	log.Fatal(http.ListenAndServe(":8001", router))
}

/** Models */

/** Handlers */
func pingAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Hello from Go api!")
}
