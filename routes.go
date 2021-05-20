package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const OneHundredMegabytes = 1 << 20 * 100

type UploadFileResponse struct {
	Key string `json:"key"`
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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
		// UploadFile()
	} else {
		// MultipartUploadFile()
	}

	// Write out the response
	res := &UploadFileResponse{Key: "b082723c-5c7a-4c37-b44d-027ff6ebc23a"}
	response, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get File ID
	key := r.URL.Query().Get("key")
	log.Println("Id supplied", key)
	_, _ = w.Write([]byte("Download Called"))
}
