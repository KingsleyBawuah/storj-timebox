package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"storj.io/uplink"
)

type server struct {
	BucketName     string
	storageProject *uplink.Project
	Router         *chi.Mux
	DB             *dynamodb.DynamoDB
}

func initBucketStorage(ctx context.Context, accessGrant, bucketName string) *uplink.Project {
	// Parse the Access Grant.
	access, err := uplink.ParseAccess(accessGrant)
	if err != nil {
		log.Fatalf("could not parse access grant: %s", err)
	}

	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		log.Fatalf("could not open project: %s", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		log.Fatalf("could not ensure bucket: %s", err)
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

func initService(ctx context.Context) *server {
	return &server{
		BucketName:     envMust("STORJ_BUCKET_NAME"),
		storageProject: initBucketStorage(ctx, envMust("STORJ_ACCESS_GRANT"), envMust("STORJ_BUCKET_NAME")),
		Router:         chi.NewRouter(),
		DB:             initDynamoDB(envMust("DYNAMO_DB_AWS_REGION"), envMust("DYNAMO_DB_ENDPOINT")),
	}
}

func main() {
	ctx := context.Background()
	// Set up http router and logger middleware.
	service := initService(ctx)
	service.Router.Use(middleware.Logger)
	service.routes()
	service.ensureTables(envMust("DYNAMO_DB_TABLE_NAME"))

	port := envMust("PORT")

	log.Println("Listening on port", port)

	// Listen for requests.
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), service.Router)
}
