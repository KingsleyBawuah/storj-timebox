package main

import (
	"log"
	"net/http"

	"github.com/KingsleyBawuah/storj-timebox/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// TODO: Set env vars for storj api keys.

func main() {
	// Set up http router and logger middleware.
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define routes.
	r.Get("/v1/file", api.DownloadFileHandler)
	r.Post("/v1/file", api.UploadFileHandler)

	log.Println("Listening on port 1337")

	// Listen for requests.
	_ = http.ListenAndServe(":1337", r)
}
