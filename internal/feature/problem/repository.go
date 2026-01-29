package problem

import (
	"context"
	"database/sql"
)

type ProblemRepository interface {
	Add(ctx context.Context, problem *Problem) (int, error)
	GetAll(ctx context.Context) ([]Problem, error)
	GetByID(ctx context.Context, id int) (Problem, error)
	Update(ctx context.Context, problem *Problem) (int, error)
	Delete(ctx context.Context, id int) (int, error)
}

type problemRepository struct {
	db *sql.DB
}

func NewProblemRepository(db *sql.DB) ProblemRepository {
	return &problemRepository{db: db}
}
