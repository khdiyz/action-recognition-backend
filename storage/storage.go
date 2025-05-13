package storage

import (
	"action-detector-backend/config"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage defines the MinIO client and bucket name.
type Storage struct {
	client     *minio.Client
	bucketName string
}

// NewStorage initializes a new MinIO client and sets the bucket name.
func NewStorage(cfg *config.Config) (*Storage, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	// Create bucket if it doesn't exist
	exists, errBucketExists := client.BucketExists(context.Background(), cfg.MinioBucketName)
	if errBucketExists != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %v", errBucketExists)
	}
	if !exists {
		errCreate := client.MakeBucket(context.Background(), cfg.MinioBucketName, minio.MakeBucketOptions{})
		if errCreate != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", errCreate)
		}
	}

	return &Storage{
		client:     client,
		bucketName: cfg.MinioBucketName,
	}, nil
}

// UploadFile uploads a file to the specified bucket.
func (s *Storage) UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// DownloadFile downloads a file from the specified bucket.
func (s *Storage) DownloadFile(ctx context.Context, objectName, destinationPath string) error {
	err := s.client.FGetObject(ctx, s.bucketName, objectName, destinationPath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	log.Printf("File %s downloaded successfully.", objectName)
	return nil
}

// DeleteFile deletes a file from the specified bucket.
func (s *Storage) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	log.Printf("File %s deleted successfully.", objectName)
	return nil
}

func (s *Storage) FPutObject(ctx context.Context, objectName, filePath, contentType string) error {
	_, err := s.client.FPutObject(ctx, s.bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed upload object: %v", err)
	}

	return nil
}
