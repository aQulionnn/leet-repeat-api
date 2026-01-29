package problem_list

import (
	"context"
	"database/sql"
)

type ProblemListRepository interface {
	Add(ctx context.Context, list *ProblemList) (int, error)
	GetAll(ctx context.Context) ([]ProblemList, error)
	Delete(ctx context.Context, id int) (int, error)
}

type problemListRepository struct {
	db *sql.DB
}

func NewProblemListRepository(db *sql.DB) ProblemListRepository {
	return &problemListRepository{db: db}
}
