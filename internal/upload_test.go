package internal

import (
	"context"
	"mime/multipart"
	"testing"

	"storj.io/uplink"
)

func TestMultipartUploadFile(t *testing.T) {
	type args struct {
		file multipart.File
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestUploadFile(t *testing.T) {
	type args struct {
		ctx        context.Context
		sp         *uplink.Project
		file       multipart.File
		fileKey    string
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UploadFile(tt.args.ctx, tt.args.sp, tt.args.file, tt.args.fileKey, tt.args.bucketName); (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
