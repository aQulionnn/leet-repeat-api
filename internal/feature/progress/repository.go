package progress

import (
	"context"
	"database/sql"
)

type ProgressRepository interface {
	GetByProblemAndList(ctx context.Context, problemID, problemListID int) (Progress, error)
	Upsert(ctx context.Context, progress *Progress) (int, error)
	GetDueForReview(ctx context.Context) ([]Progress, error)
}

type progressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) ProgressRepository {
	return &progressRepository{db: db}
}
