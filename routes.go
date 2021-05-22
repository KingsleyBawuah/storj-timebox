package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KingsleyBawuah/storj-timebox/internal"
	"github.com/go-chi/chi/v5"
)

const (
	OneHundredMegabytes = 1 << 20 * 100
)

func (s *server) routes() {
	s.Router.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	s.Router.HandleFunc("/file/{key}", s.DownloadFileHandler())
	s.Router.HandleFunc("/file", s.UploadFileHandler())
}

func (s *server) UploadFileHandler() http.HandlerFunc {
	type UploadFileResponse struct {
		Key string `json:"key"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Should I be doing this?
		ctx := context.Background()
		// Assert request body exists
		if r.Body == nil {
			http.Error(w, "Error request body required", http.StatusBadRequest)
		}

		// Only allow 32.5mb as the max post body size.
		r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512)

		// Read the expected data about the file.
		name := r.FormValue("name")
		maxAllowedDownloads := r.FormValue("maxAllowedDownloads")
		expirationDateTimeStr := r.FormValue("expirationDateTime")

		// Assert body fields exist.
		if name == "" || maxAllowedDownloads == "" || expirationDateTimeStr == "" {
			http.Error(w, "Error, all request body fields required", http.StatusBadRequest)
		}

		// Parse the time to validate it.
		_, err := time.Parse(time.RFC3339, expirationDateTimeStr)
		if err != nil {
			http.Error(w, "Error parsing expirationDateTime, please use RFC3339 date-time", http.StatusBadRequest)
		}

		// Read the file
		f, fh, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer f.Close()

		// Determine which upload method is appropriate and begin uploading the file to the Storj DCS network.
		if fh.Size < OneHundredMegabytes {
			if err := internal.UploadFile(ctx, s.storageProject, f, fh.Filename, s.BucketName); err != nil {
				// TODO: Populate error message.
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			// MultipartUploadFile()
		}

		// Write out the response
		res := &UploadFileResponse{Key: fh.Filename}
		response, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(response)
	}
}

func (s *server) DownloadFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Can/Should I avoid declaring ctx in each handler?
		ctx := context.Background()
		// Get File ID
		key := chi.URLParam(r, "key")
		log.Printf("Download of file %s requested \n", key)

		file, err := internal.DownloadFile(ctx, s.storageProject, key, s.BucketName)
		if err != nil {
			log.Println("error downloading file", err.Error())
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", key))
		_, _ = w.Write(file)
	}
}
