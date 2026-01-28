package upload

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

// StorageService is a Firebase Storage implementation of the upload service.
type StorageService struct {
	client     *storage.Client
	bucketName string
}

// NewStorageService creates a new Storage-backed upload service.
func NewStorageService(client *storage.Client, bucketName string) *StorageService {
	return &StorageService{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *StorageService) UploadImage(ctx context.Context, fileData []byte, contentType string, uploadType string) (UploadResult, error) {
	// Generate a unique ID for the file
	id := uuid.New().String()

	// Determine file extension based on content type
	ext := getFileExtension(contentType)

	// Create the object path: type/id.ext
	objectPath := fmt.Sprintf("%s/%s%s", uploadType, id, ext)

	// Get bucket handle
	bucket := s.client.Bucket(s.bucketName)

	// Create object handle
	obj := bucket.Object(objectPath)

	// Create a writer for the object
	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	// Write the file data
	if _, err := writer.Write(fileData); err != nil {
		log.Printf("Failed to write file data: %v", err)
		return UploadResult{}, fmt.Errorf("failed to write file data: %w", err)
	}

	// Close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		log.Printf("Failed to close writer: %v", err)
		return UploadResult{}, fmt.Errorf("failed to close writer: %w", err)
	}

	// Make the object publicly readable
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Printf("Failed to set ACL: %v", err)
		return UploadResult{}, fmt.Errorf("failed to set ACL: %w", err)
	}

	// Generate the public URL with resized image suffix (Firebase Resize Image Extension)
	resizedObjectPath := fmt.Sprintf("%s/%s_200x200%s", uploadType, id, ext)
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, resizedObjectPath)

	return UploadResult{
		ID:   id,
		URL:  url,
		Type: uploadType,
	}, nil
}

// GetSignedURL generates a signed URL for the object (optional, for private access)
func (s *StorageService) GetSignedURL(ctx context.Context, objectPath string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := s.client.Bucket(s.bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// getFileExtension returns the file extension based on content type
func getFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ".jpg"
	}
}
