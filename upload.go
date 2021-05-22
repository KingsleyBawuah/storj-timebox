package main

import (
	"mime/multipart"

	_ "storj.io/uplink"
)

// Uploads a file, this is used for smaller files.
func UploadFile(file multipart.File) {
}

// Uploads a single file, this is for larger files that need to be uploaded in a more efficient fashion.
func MultipartUploadFile(file multipart.File) {
}
