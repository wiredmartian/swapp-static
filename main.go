package main

import (
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
type Book struct {
	ID     string
	Isbn   string
	Author *Author
}

type Author struct {
	name string
}

/** Handlers */
/** Looks similar to nodejs pingAPI(req, res) */
func pingAPI(w http.ResponseWriter, r *http.Request) {

}
