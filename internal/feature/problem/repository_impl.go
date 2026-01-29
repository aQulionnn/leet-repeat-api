package problem

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *problemRepository) Add(ctx context.Context, problem *Problem) (int, error) {
	query := `INSERT INTO problem (question, difficulty) VALUES ($1, $2) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, problem.Question, problem.Difficulty).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add problem: %w", err)
	}

	return id, nil
}

func (r *problemRepository) GetAll(ctx context.Context) ([]Problem, error) {
	query := `SELECT id, question, difficulty FROM problem ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all problems: %w", err)
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(&problem.ID, &problem.Question, &problem.Difficulty); err != nil {
			return nil, fmt.Errorf("failed to scan problem: %w", err)
		}

		problems = append(problems, problem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating problems: %w", err)
	}

	return problems, nil
}

func (r *problemRepository) GetByID(ctx context.Context, id int) (Problem, error) {
	query := `SELECT id, question, difficulty FROM problem WHERE id = $1`

	var problem Problem
	err := r.db.QueryRowContext(ctx, query, id).Scan(&problem.ID, &problem.Question, &problem.Difficulty)

	if err == sql.ErrNoRows {
		return Problem{}, fmt.Errorf("problem not found with id %d", id)
	}
	if err != nil {
		return Problem{}, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (r *problemRepository) Update(ctx context.Context, problem *Problem) (int, error) {
	query := `UPDATE problem SET question = $1, difficulty = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, problem.Question, problem.Difficulty, problem.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to update problem: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func (r *problemRepository) Delete(ctx context.Context, id int) (int, error) {
	query := `DELETE FROM problem WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete problem: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}
