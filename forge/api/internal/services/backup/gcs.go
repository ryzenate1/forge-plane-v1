package backup

import (
	"context"
	"fmt"
)

type GCSAdapter struct {
	bucketName string
	keyFile    string
}

func NewGCSAdapter(bucketName, keyFile string) *GCSAdapter {
	return &GCSAdapter{
		bucketName: bucketName,
		keyFile:    keyFile,
	}
}

func (a *GCSAdapter) Name() string { return "gcs" }

func (a *GCSAdapter) Upload(ctx context.Context, path string, data []byte) error {
	return nil
}

func (a *GCSAdapter) Download(ctx context.Context, path string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *GCSAdapter) Delete(ctx context.Context, path string) error {
	return nil
}

func (a *GCSAdapter) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}

func (a *GCSAdapter) Exists(ctx context.Context, path string) (bool, error) {
	return false, nil
}
