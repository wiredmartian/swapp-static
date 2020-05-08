package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	/** Router handlers */
	router.HandleFunc("/api/ping", pingAPI).Methods("GET")
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":8001", router))
}

/** Models */

/** Handlers */
func pingAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Hello from Go api!")
}
func getFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/11290243.jpeg")
}
