package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.Use(requestLogging)

	/** Router handlers */
	router.HandleFunc("/api/ping", pingAPI).Methods("GET")
	router.HandleFunc("/api/upload", uploadFileHandler).Methods("POST")
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":8001", router))
	fmt.Println("Server running on PORT 8001")
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
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 * 1024 * 1024)
	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer file.Close()
	fmt.Printf("%v\n", header.Filename)
	fmt.Printf("%v\n", header.Size)
	fmt.Printf("%v\n", header.Header)

	tempFile, error := ioutil.TempFile("static", "upload-*.png")
	if error != nil {
		fmt.Println(error)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(error)
		return
	}
	defer tempFile.Close()

	fileBytes, er := ioutil.ReadAll(file)
	if er != nil {
		fmt.Println(err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}
	tempFile.Write(fileBytes)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode("Successfully uploaded file")
}

/** Middleware */
func requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RequestURI)
		next.ServeHTTP(writer, request)
	})
}
