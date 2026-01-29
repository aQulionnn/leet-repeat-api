package problem_list_problem

import (
	"context"
	"fmt"
)

func (r *problemListProblemRepository) Add(ctx context.Context, plp *ProblemListProblem) (int, error) {
	query := `INSERT INTO problem_list_problem (problem_id, problem_list_id) VALUES ($1, $2) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, plp.ProblemID, plp.ProblemListID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add problem to list: %w", err)
	}

	return id, nil
}

func (r *problemListProblemRepository) Remove(ctx context.Context, problemID, problemListID int) (int, error) {
	query := `DELETE FROM problem_list_problem WHERE problem_id = $1 AND problem_list_id = $2`

	result, err := r.db.ExecContext(ctx, query, problemID, problemListID)
	if err != nil {
		return 0, fmt.Errorf("failed to remove problem from list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func (r *problemListProblemRepository) GetByList(ctx context.Context, problemListID int) ([]ProblemListProblem, error) {
	query := `SELECT id, problem_id, problem_list_id FROM problem_list_problem WHERE problem_list_id = $1 ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, problemListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problems by list: %w", err)
	}
	defer rows.Close()

	var items []ProblemListProblem
	for rows.Next() {
		var item ProblemListProblem
		if err := rows.Scan(&item.ID, &item.ProblemID, &item.ProblemListID); err != nil {
			return nil, fmt.Errorf("failed to scan problem list problem: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating problem list problems: %w", err)
	}

	return items, nil
}
