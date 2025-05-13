package usecase

import (
	"action-detector-backend/config"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/pkg/response"
	"action-detector-backend/storage"
	"context"
	"io"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type fileUsecase struct {
	storage *storage.Storage
	logger  *logger.Logger
	cfg     *config.Config
}

func newFileUsecase(storage *storage.Storage, logger *logger.Logger, cfg *config.Config) *fileUsecase {
	return &fileUsecase{
		storage: storage,
		logger:  logger,
		cfg:     cfg,
	}
}

func (s *fileUsecase) UploadFile(ctx context.Context, file io.Reader, fileSize int64, contentType string) (string, error) {
	// Generate a unique object name for the file
	objectName := uuid.NewString()

	// Upload the file directly using the storage service
	err := s.storage.UploadFile(ctx, objectName, file, fileSize, contentType)
	if err != nil {
		return "", response.ServiceError(err, codes.Internal)
	}

	fileLink := "https://" + s.cfg.MinioEndpoint + "/" + s.cfg.MinioBucketName + "/" + objectName

	// Return the object name (file ID)
	return fileLink, nil
}
