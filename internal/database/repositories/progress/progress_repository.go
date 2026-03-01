package progress

import (
	"context"
	"database/sql"
	"leet-repeat-api/internal/database/models"
)

type ProgressRepository interface {
	BulkUpsert(ctx context.Context, progressList *[]models.Progress) (int, error)
	GetAll(ctx context.Context) ([]models.Progress, error)
	Clear(ctx context.Context) (int, error)
}

type progressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) ProgressRepository {
	return &progressRepository{db: db}
}
