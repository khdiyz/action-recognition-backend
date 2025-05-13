package usecase

import (
	"action-detector-backend/models"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/pkg/response"
	"action-detector-backend/repository"
	"context"

	"google.golang.org/grpc/codes"
)

type actionUsecase struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func newActionUsecase(repo *repository.Repository, logger *logger.Logger) *actionUsecase {
	return &actionUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *actionUsecase) CreateAction(ctx context.Context, action models.Action) error {
	err := u.repo.ActionRepo.CreateAction(ctx, action)
	if err != nil {
		return response.ServiceError(err, codes.Internal)
	}

	return nil
}

func (u *actionUsecase) GetActions(ctx context.Context) ([]models.Action, error) {
	actions, err := u.repo.ActionRepo.GetActions(ctx)
	if err != nil {
		return nil, response.ServiceError(err, codes.Internal)
	}

	return actions, nil
}

func (u *actionUsecase) DeleteActions(ctx context.Context) error {
	err := u.repo.ActionRepo.DeleteActions(ctx)
	if err != nil {
		return response.ServiceError(err, codes.Internal)
	}

	return nil
}
