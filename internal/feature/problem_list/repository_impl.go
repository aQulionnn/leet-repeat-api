package problem_list

import (
	"context"
	"fmt"
)

func (r *problemListRepository) Add(ctx context.Context, list *ProblemList) (int, error) {
	query := `INSERT INTO problem_list (name) VALUES ($1) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, list.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add problem list: %w", err)
	}

	return id, nil
}

func (r *problemListRepository) GetAll(ctx context.Context) ([]ProblemList, error) {
	query := `SELECT id, name FROM problem_list ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all problem lists: %w", err)
	}
	defer rows.Close()

	var lists []ProblemList
	for rows.Next() {
		var list ProblemList
		if err := rows.Scan(&list.ID, &list.Name); err != nil {
			return nil, fmt.Errorf("failed to scan problem list: %w", err)
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating problem lists: %w", err)
	}

	return lists, nil
}

func (r *problemListRepository) Delete(ctx context.Context, id int) (int, error) {
	query := `DELETE FROM problem_list WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete problem list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}
