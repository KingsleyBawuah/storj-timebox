package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"storj.io/uplink"
)

var sp *uplink.Project

var bucketName string

func initBucketStorage(ctx context.Context, accessGrant, bucketName string) *uplink.Project {
	// Parse the Access Grant.
	access, err := uplink.ParseAccess(accessGrant)
	if err != nil {
		log.Fatalf("could not parse access grant: %s", err.Error())
	}

	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		log.Fatalf("could not open project: %s", err.Error())
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		log.Fatalf("could not ensure bucket: %s", err.Error())
	}

	return project
}

// fetches and returns the given env variable, fatals and
// captures an exception if the variable is an empty string
func envMust(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		e := errors.New("environment variable missing - " + varName)
		log.Fatalln(e.Error())
	}
	return value
}

func main() {
	// TODO: What else can I do with this to help solve this problem?
	ctx := context.Background()

	bucketName = envMust("STORJ_BUCKET_NAME")

	sp = initBucketStorage(ctx, envMust("STORJ_ACCESS_GRANT"), envMust("STORJ_BUCKET_NAME"))

	// Set up http router and logger middleware.
	r := chi.NewRouter()
	// TODO: Handle this manually with a structured logging library.
	r.Use(middleware.Logger)

	// TODO: Handle routing better.
	// Define routes with additional health check for monitoring.
	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	r.Get("/file/{key}", DownloadFileHandler)
	r.Post("/file", UploadFileHandler)

	port := envMust("PORT")

	log.Println("Listening on port", port)

	// Listen for requests.
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
