package problem_list_problem

import (
	"context"
	"database/sql"
)

type ProblemListProblemRepository interface {
	Add(ctx context.Context, plp *ProblemListProblem) (int, error)
	Remove(ctx context.Context, problemID, problemListID int) (int, error)
	GetByList(ctx context.Context, problemListID int) ([]ProblemListProblem, error)
}

type problemListProblemRepository struct {
	db *sql.DB
}

func NewProblemListProblemRepository(db *sql.DB) ProblemListProblemRepository {
	return &problemListProblemRepository{db: db}
}
