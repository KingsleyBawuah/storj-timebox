package internal

import (
	"context"
	"io/ioutil"

	"storj.io/uplink"
)

// Downloads a file using Storj DCS
func DownloadFile(ctx context.Context, sp *uplink.Project, key, bucketName string) ([]byte, error) {
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
}
