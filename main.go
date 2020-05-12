package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.Use(requestLogging)
	router.Use(parseToken)

	/** Router handlers */
	router.HandleFunc("/api/ping", pingAPI).Methods("GET")
	router.HandleFunc("/api/upload", uploadFileHandler).Methods("POST")
	router.HandleFunc("/static/{filename}", getFile).Methods("GET")
	fmt.Println("Server running on PORT 8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

/** Models */
type FileUploadRes struct {
	Message string
	FileUrl string
}

/** Handlers */
func pingAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode("Hello from Go api!")
}
func getFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]
	http.ServeFile(w, r, "./static/"+filename)
}
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(5 * 1024 * 1024)
	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	defer file.Close()
	fmt.Printf("%v\n", header.Filename)
	fmt.Printf("%v\n", header.Size)
	fmt.Printf("%v\n", header.Header)

	tempFile, _error := ioutil.TempFile("static", "upload-*.png")
	if _error != nil {
		fmt.Println(_error)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(_error)
		return
	}
	defer tempFile.Close()

	fileBytes, er := ioutil.ReadAll(file)
	if er != nil {
		fmt.Println(err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	bytesW, _err := tempFile.Write(fileBytes)
	if _err != nil {
		fmt.Println(err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(_err)
		return
	}
	fmt.Println(bytesW)
	res := FileUploadRes{Message: "Successfully uploaded file", FileUrl: tempFile.Name()}
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func parseToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		requestToken := authHeader[7:len(authHeader)]
		fmt.Println(requestToken)
		tokenString := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI5YjEwYzg5OC0wMTI4LTRlMDQtODNiOC05MzIxMjNjYzM3ZDMiLCJ1c2VyRnVsbE5hbWUiOiI5YjEwYzg5OC0wMTI4LTRlMDQtODNiOC05MzIxMjNjYzM3ZDMiLCJyZXNlbGxlcklkIjpudWxsLCJzaXRlS2V5IjoicG9ydGFsLm15dGVsbmV0LmNvLnphIiwiYWRkaXQiOnt9LCJpYXQiOjE1ODkzMTc1NzcsImV4cCI6MTU4OTMyMTE3NywiYXVkIjoicG9ydGFsLm15dGVsbmV0LmNvLnphIiwiaXNzIjoic3VwcG9ydEBteXRlbG5ldC5jby56YSIsInN1YiI6IjliMTBjODk4LTAxMjgtNGUwNC04M2I4LTkzMjEyM2NjMzdkMyJ9.4WI5XjOPH6Xy36ykjOeB26XnqKGPmgTF_-Vd3eYAWA1BPUNtX9RpZTlDLSQw09H-UwzKrhGj2d2RUTEHXDONpg"
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("8D6049A45555471584B0CADC2E2B8A45"), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(token.Claims)
			fmt.Println(claims["foo"], claims["nbf"])
			next.ServeHTTP(writer, request)
		} else {
			fmt.Println(err)
		}
	})
}

/** Middleware */
func requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RequestURI)
		next.ServeHTTP(writer, request)
	})
}
