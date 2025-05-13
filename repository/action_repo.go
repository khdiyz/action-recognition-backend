package repository

import (
	"action-detector-backend/models"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/pkg/postgres"
	"context"
)

type actionRepo struct {
	db     *postgres.Postgres
	logger *logger.Logger
}

func newActionRepo(db *postgres.Postgres, logger *logger.Logger) *actionRepo {
	return &actionRepo{
		db:     db,
		logger: logger,
	}
}

func (r *actionRepo) CreateAction(ctx context.Context, action models.Action) error {
	query := `
	insert into actions (
		video_url, predicted_actions
	) values ($1, $2);`

	_, err := r.db.Pool.Exec(ctx, query, action.VideoURL, action.PredictedActions)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *actionRepo) GetActions(ctx context.Context) ([]models.Action, error) {
	actions := []models.Action{}

	query := `
	select id, predicted_actions, video_url, created_at
	from actions`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var action models.Action

		if err = rows.Scan(
			&action.Id,
			&action.PredictedActions,
			&action.VideoURL,
			&action.CreatedAt,
		); err != nil {
			r.logger.Error(err)
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (r *actionRepo) DeleteActions(ctx context.Context) error {
	query := `delete from actions`

	_, err := r.db.Pool.Exec(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}
