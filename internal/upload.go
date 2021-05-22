package internal

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"storj.io/uplink"
)

// Uploads a file, this is used for smaller files.
func UploadFile(ctx context.Context, sp *uplink.Project, file multipart.File, fileKey, bucketName string) error {
	// TODO: Set object info so the file expires based on the time passed in and custom metadata containing the download counter is set.
	upload, err := sp.UploadObject(ctx, bucketName, fileKey, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
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

	return nil
}

// Uploads a single file, this is for larger files that need to be uploaded in a more efficient fashion.
func MultipartUploadFile(file multipart.File) {
}
