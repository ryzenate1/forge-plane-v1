package backup

import (
	"context"
	"fmt"
)

type AzureAdapter struct {
	containerName string
	connString    string
}

func NewAzureAdapter(containerName, connString string) *AzureAdapter {
	return &AzureAdapter{
		containerName: containerName,
		connString:    connString,
	}
}

func (a *AzureAdapter) Name() string { return "azure" }

func (a *AzureAdapter) Upload(ctx context.Context, path string, data []byte) error {
	return nil
}

func (a *AzureAdapter) Download(ctx context.Context, path string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *AzureAdapter) Delete(ctx context.Context, path string) error {
	return nil
}

func (a *AzureAdapter) List(ctx context.Context, prefix string) ([]string, error) {
	return nil, nil
}

func (a *AzureAdapter) Exists(ctx context.Context, path string) (bool, error) {
	return false, nil
}
