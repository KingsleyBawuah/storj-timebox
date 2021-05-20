package main

import (
	"log"
	"mime/multipart"
	"net/http"
	// TODO: Set env vars for storj api keys.
	_ "storj.io/uplink"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Uploads a file, this is used for smaller files.
func UploadFile(file multipart.File) {
}

// Uploads a single file, this is for larger files that need to be uploaded in a more efficient fashion.
func MultipartUploadFile(file multipart.File) {
}

// Downloads a file.
func DownloadFile() []byte {
	return nil
}

func main() {
	// Set up http router and logger middleware.
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define routes.
	r.Get("/v1/file", DownloadFileHandler)
	r.Post("/v1/file", UploadFileHandler)

	log.Println("Listening on port 1337")

	// Listen for requests.
	_ = http.ListenAndServe(":1337", r)
}
