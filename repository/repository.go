package repository

import (
	"action-detector-backend/models"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/pkg/postgres"
	"context"
)

type Repository struct {
	ActionRepo Action
}

func NewRepository(db *postgres.Postgres, logger *logger.Logger) *Repository {
	return &Repository{
		ActionRepo: newActionRepo(db, logger),
	}
}

type Action interface {
	CreateAction(ctx context.Context, action models.Action) error
	GetActions(ctx context.Context) ([]models.Action, error)
	DeleteActions(ctx context.Context) error
}
