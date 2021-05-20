package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "storj.io/uplink"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UploadFileResponse struct {
	Id string `json:"id"`
}

// TODO: Break these methods out so all the request/body logic is separate from the actual "business logic"
func UploadFile(w http.ResponseWriter, r *http.Request) {
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
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()

	// Write out the response
	res := &UploadFileResponse{Id: "b082723c-5c7a-4c37-b44d-027ff6ebc23a"}
	response, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}

// TODO: Break these methods out so all the request/body logic is separate from the actual "business logic"
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Get File ID
	id := r.URL.Query().Get("id")
	log.Println("Id supplied", id)
	_, _ = w.Write([]byte("Download Called"))
}

func main() {
	// Set up http router and logger middleware.
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define routes.
	r.Get("/v1/file", DownloadFile)
	r.Post("/v1/file", UploadFile)

	log.Println("Listening on port 1337")

	// Listen for requests.
	_ = http.ListenAndServe(":1337", r)
}
