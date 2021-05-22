package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KingsleyBawuah/storj-timebox/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"storj.io/uplink"
)

type storjProject *uplink.Project

func initBucketStorage(ctx context.Context, ag, bn string) (*uplink.Project, error) {
	// Parse the Access Grant.
	access, err := uplink.ParseAccess(ag)
	if err != nil {
		return nil, fmt.Errorf("could not parse access grant: %v", err)
	}

	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, bn)
	if err != nil {
		return nil, fmt.Errorf("could not ensure bucket: %v", err)
	}

	return project, nil
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
	// Set up http router and logger middleware.
	r := chi.NewRouter()
	// TODO: Handle this manually with a structured logging library.
	r.Use(middleware.Logger)

	// TODO: Handle routing better.
	// Define routes.
	r.Get("/v1/file", api.DownloadFileHandler)
	r.Post("/v1/file", api.UploadFileHandler)

	log.Println("Listening on port 8080")

	// Listen for requests.
	_ = http.ListenAndServe(":8080", r)
}
