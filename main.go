package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	router.Use(requestLogging)
	router.Use(parseToken)

	/** Router handlers */
	router.HandleFunc("/api/ping", pingAPI).Methods("GET")
	router.HandleFunc("/api/upload", uploadFileHandler).Methods("POST")
	router.HandleFunc("/static/{filename}", getFile).Methods("GET")

	/** load .env */
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env")
	}
	fmt.Println("Server running on PORT 8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

/** Models */
type FileUploadRes struct {
	Message string
	FileUrl string
}
type FileUploads struct {
	FileUrls []string
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
func uploadMultipleFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Uploading files...")
	_ = r.ParseMultipartForm(5 * 1024 * 1024)
	fileHeaders := r.MultipartForm.File["images"]
	var uploadedFiles []string
	for _, fh := range fileHeaders {
		file, err := fh.Open()
		if err == nil {
			fileURI, _err := uploadFile(file)
			if _err == nil {
				uploadedFiles = append(uploadedFiles, fileURI)
			}
		}
	}
	fileUrls := FileUploads{FileUrls: uploadedFiles}
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(fileUrls)
}
func uploadFile(f multipart.File) (fileURI string, error error) {
	defer f.Close()
	tempFile, _ := ioutil.TempFile("static", "upload-*.png")
	defer tempFile.Close()
	fileBytes, _ := ioutil.ReadAll(f)
	_, err := tempFile.Write(fileBytes)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil

}
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Uploading file...")
	_ = r.ParseMultipartForm(5 * 1024 * 1024)
	file, header, err := r.FormFile("image")
	fmt.Println(file)
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
		token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := os.Getenv("JWTSECRET")
			return []byte(secret), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			/** verified token here*/
			fmt.Println(token.Claims)
			fmt.Println(claims["iss"])
			next.ServeHTTP(writer, request)
		} else {
			fmt.Println(err)
			writer.Header().Set("content-type", "application/json")
			writer.WriteHeader(400)
			_ = json.NewEncoder(writer).Encode(err.Error())
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
