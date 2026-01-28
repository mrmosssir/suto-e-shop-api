package upload

import "context"

// UploadResult represents the result of an upload operation.
type UploadResult struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

// Service provides upload operations.
type Service interface {
	UploadImage(ctx context.Context, fileData []byte, contentType string, uploadType string) (UploadResult, error)
}
