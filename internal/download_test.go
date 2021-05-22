package internal

import (
	"context"
	"reflect"
	"testing"

	"storj.io/uplink"
)

func TestDownloadFile(t *testing.T) {
	type args struct {
		ctx        context.Context
		sp         *uplink.Project
		key        string
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadFile(tt.args.ctx, tt.args.sp, tt.args.key, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
