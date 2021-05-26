package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"storj.io/uplink"
)

// Uploads a file, this is used for smaller files (under 100mb)
func UploadFile(ctx context.Context, sp *uplink.Project, DB *server, file multipart.File, fileKey, bucketName string, maxDownloads int, expires time.Time) error {
	log.Println("maxDownloads", maxDownloads)
	// Prevent file name collisions.
	_, err := fetchFileObj(ctx, sp, fileKey, bucketName)
	if err != nil {
		// TODO: Not the biggest fan of relying on this string. If Storj changed the error message this would stop working
		if !strings.Contains(err.Error(), "object not found") {
			return err
		}
	}

	upload, err := sp.UploadObject(ctx, bucketName, fileKey, &uplink.UploadOptions{
		Expires: expires,
	})
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}

	err = upload.SetCustomMetadata(ctx, uplink.CustomMetadata{
		downLoadLimitKey: strconv.Itoa(maxDownloads),
	})
	if err != nil {
		return fmt.Errorf("could not set metadata on uploaded file %v", err)
	}

	// Copy the data to the upload.
	_, err = io.Copy(upload, file)
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}

	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}

	// Create the download limit counter.
	err = DB.CreateDownloadCount(fileKey, os.Getenv("DYNAMO_DB_TABLE_NAME"))
	if err != nil {
		return fmt.Errorf("could not create download count in db: %v", err)
	}

	return nil
}

// Uploads a single file, this is for larger files (over 100mb) that need to be uploaded in a more efficient fashion.
func MultipartUploadFile(file multipart.File) {
}
