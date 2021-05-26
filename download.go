package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"storj.io/uplink"
)

const (
	downLoadLimitKey = "timebox:downloadLimit"
)

// Fetch information about the file from storjDCS
func fetchFileObj(ctx context.Context, sp *uplink.Project, key, bucketName string) (*uplink.Object, error) {
	obj, err := sp.StatObject(ctx, bucketName, key)
	if err != nil {
		return nil, fmt.Errorf("there is an issue fetching the file %v", err)
	}

	return obj, nil
}

// Validate that we haven't reached the max download count.
func validateDownload(fileObj *uplink.Object, downloadCount int) bool {
	downloadLimit := fileObj.Custom[downLoadLimitKey]

	// Prevent files without a downloadLimit from being downloaded.
	if downloadLimit == "" {
		return false
	}

	dl, err := strconv.Atoi(downloadLimit)
	if err != nil {
		log.Println(err)
		return false
	}

	if downloadCount >= dl {
		return false
	}

	return true
}

// Downloads a file using Storj DCS
func DownloadFile(ctx context.Context, sp *uplink.Project, s *server, key, bucketName string) ([]byte, error) {
	obj, err := fetchFileObj(ctx, sp, key, bucketName)
	if err != nil {
		return nil, err
	}

	// TODO: Remove the if, send that we can't find the file somehow.
	if obj != nil {
		log.Printf("File exists!!!!, here's some metadata %+v\\n", obj.Custom)

		count, err := s.GetDownloadCount(key, os.Getenv("DYNAMO_DB_TABLE_NAME"))
		if err != nil {
			return nil, err
		}

		// TODO: Read files row in dynamodb for count.
		if validateDownload(obj, count) {
			// TODO: Look into if there is an extra benefit that the extra download options unlock, especially for larger files.
			download, err := sp.DownloadObject(ctx, bucketName, key, nil)
			if err != nil {
				return nil, err
			}
			// Read everything from the download stream
			// TODO: Don't read the entire file into memory like this.
			receivedContents, err := ioutil.ReadAll(download)
			if err != nil {
				return nil, err
			}
			return receivedContents, nil
		} else {
			return nil, errors.New("can't download file, limit reached")
		}
	} else {
		return nil, errors.New("file not found")
	}
}
