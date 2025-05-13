package usecase

import (
	"action-detector-backend/config"
	"action-detector-backend/models"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/repository"
	"action-detector-backend/storage"
	"context"
	"io"
)

type Usecase struct {
	Action
	File
}

type Dependencies struct {
	Repository *repository.Repository
	Logger     *logger.Logger
	Config     *config.Config
	Storage    *storage.Storage
}

func NewUsecase(deps Dependencies) *Usecase {
	return &Usecase{
		Action: newActionUsecase(deps.Repository, deps.Logger),
		File:   newFileUsecase(deps.Storage, deps.Logger, deps.Config),
	}
}

type Action interface {
	CreateAction(ctx context.Context, action models.Action) error
	GetActions(ctx context.Context) ([]models.Action, error)
	DeleteActions(ctx context.Context) error
}

type File interface {
	UploadFile(ctx context.Context, file io.Reader, fileSize int64, contentType string) (string, error)
}
