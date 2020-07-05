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
	"path/filepath"
	"strings"
)

func main() {
	router := mux.NewRouter()
	router.Use(requestLogging)
	router.Use(parseToken)

	/** Router handlers */
	router.HandleFunc("/api/health", health).Methods("GET")
	router.HandleFunc("/api/upload", uploadFilesHandler).Methods("POST")
	router.HandleFunc("/api/create-dir", createDirHandler).Methods("POST")
	/** I need to post file paths */
	router.HandleFunc("/api/purge", purgeDirsHandler).Methods("POST")
	router.HandleFunc("/static/{dir}/{filename}", getFileHandler).Methods("GET")
	router.HandleFunc("/static/{filename}", getFileHandler).Methods("GET")

	/** load .env */
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env")
	}
	fmt.Println("Server running on PORT 8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

/** Models */
type FileUploads struct {
	Message  string
	FileUrls []string
}
type Purge struct {
	FileDir string
}
type CreateDir struct {
	DirName string
}
type FolderStructure struct {
	FileDir string
	Files   []os.FileInfo
}
type MessageResponse struct {
	Message string
}

/** Handlers */
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(`{"alive"}: true`)
}

func uploadFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Uploading files...")
	_ = r.ParseMultipartForm(5 * 1024 * 1024)
	fileHeaders := r.MultipartForm.File["images"]
	dir := r.MultipartForm.Value["dir"][0]
	var uploadedFiles []string
	for _, fh := range fileHeaders {
		file, err := fh.Open()
		if err == nil {
			fileURI, _err := uploadFile(file, dir)
			if _err == nil {
				uploadedFiles = append(uploadedFiles, fileURI)
			}
		}
	}
	resMessage := fmt.Sprintf("%v out of %v files were uploaded", len(uploadedFiles), len(fileHeaders))
	fileUrls := FileUploads{Message: resMessage, FileUrls: uploadedFiles}
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(fileUrls)
}

/*
Removes all files that exist on a dir
*/
func purgeDirsHandler(w http.ResponseWriter, r *http.Request) {
	var body Purge
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println(err)
	}
	var filePaths []string = strings.Split(body.FileDir, ",")
	counter := 0
	for _, path := range filePaths {
		/** we don't wanna remove a wrong dir */
		path = filepath.Join("./static/", path)
		if strings.Index("/static/profiles", path) >= 0 || strings.Index("/static/products", path) >= 0 {
			err := os.RemoveAll(path)
			if err == nil {
				counter++
			}
			fmt.Println(err)
		}
	}
	responseMessage := fmt.Sprintf("Remove %v folders with file", counter)
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(responseMessage)
}
func getFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]
	dir := params["dir"]

	if dir != "" {
		fileInfo, err := os.Stat("./static/" + dir + "/" + filename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "./static/"+dir+"/"+fileInfo.Name())
	} else {
		_fileInfo, err := os.Stat("./static/" + filename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "./static/"+_fileInfo.Name())
	}
}

func createDirHandler(w http.ResponseWriter, r *http.Request) {
	var res MessageResponse
	var dirName CreateDir
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dirName)
	if err != nil {
		fmt.Println(err)
	}
	dir := strings.TrimSpace(dirName.DirName)
	w.Header().Set("content-type", "application/json")
	if strings.ContainsAny(dir, "<>[]{}()!@#$^&*=+;:|/?.,`~") || strings.Contains(dir, " ") {
		w.WriteHeader(http.StatusBadRequest)
		res.Message = fmt.Sprintf("%v contains empty space and/or special characters", dir)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	dirInfo, _err := createDir(dir)
	if _err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Message = _err.Error()
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	w.WriteHeader(http.StatusOK)
	res.Message = fmt.Sprintf("%v folder was successfully created", dirInfo.Name())
	_ = json.NewEncoder(w).Encode(res)
}

/** END Handlers */

func createDir(dir string) (dirInfo os.FileInfo, error error) {
	newPath := filepath.Join("./static", dir)
	fileInfo, err := os.Stat(newPath)
	/** if the dir does not exist, create it*/
	if os.IsNotExist(err) && fileInfo == nil {
		_err := os.Mkdir(newPath, os.ModePerm)
		if _err != nil {
			return nil, _err
		}
		newDirInfo, _ := os.Stat(newPath)
		return newDirInfo, nil
	}
	return nil, os.ErrExist
}
func uploadFile(f multipart.File, dir string) (fileURI string, error error) {
	defer f.Close()
	newPath := filepath.Join("./static", dir)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		_err := os.MkdirAll(newPath, os.ModePerm)
		if _err != nil {
			return "", _err
		}
	}
	tempFile, er := ioutil.TempFile(newPath, "upload-*.png")
	if er != nil {
		return "", er
	}
	defer tempFile.Close()
	fileBytes, _ := ioutil.ReadAll(f)
	_, err := tempFile.Write(fileBytes)
	if err != nil {
		return "", err
	}
	return "/" + tempFile.Name(), nil

}

/** Middleware */
func parseToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" && request.Method == "POST" {
			writer.Header().Set("content-type", "application/json")
			writer.WriteHeader(401)
			_ = json.NewEncoder(writer).Encode(`{"message": "Authorization token not found"}`)
			return
		}
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
			fmt.Println(claims["userId"])
			next.ServeHTTP(writer, request)
		} else {
			writer.Header().Set("content-type", "application/json")
			writer.WriteHeader(400)
			_ = json.NewEncoder(writer).Encode(err.Error())
		}
	})
}

func requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.RequestURI)
		next.ServeHTTP(writer, request)
	})
}
