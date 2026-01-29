package progress

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (r *progressRepository) GetByProblemAndList(ctx context.Context, problemID, problemListID int) (Progress, error) {
	query := `SELECT id, perceived_difficulty, status, last_solved_at, next_review_at, problem_id, problem_list_id 
	          FROM progress 
	          WHERE problem_id = $1 AND problem_list_id = $2`

	var progress Progress
	err := r.db.QueryRowContext(ctx, query, problemID, problemListID).Scan(
		&progress.ID,
		&progress.PerceivedDifficulty,
		&progress.Status,
		&progress.LastSolvedAt,
		&progress.NextReviewAt,
		&progress.ProblemID,
		&progress.ProblemListID,
	)

	if err == sql.ErrNoRows {
		return Progress{}, fmt.Errorf("progress not found for problem %d in list %d", problemID, problemListID)
	}
	if err != nil {
		return Progress{}, fmt.Errorf("failed to get progress: %w", err)
	}

	return progress, nil
}

func (r *progressRepository) Upsert(ctx context.Context, progress *Progress) (int, error) {
	query := `UPDATE progress 
	          SET perceived_difficulty = $1, 
	              status = $2, 
	              last_solved_at = $3, 
	              next_review_at = $4
	          WHERE problem_id = $5 AND problem_list_id = $6
	          RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		progress.PerceivedDifficulty,
		progress.Status,
		progress.LastSolvedAt,
		progress.NextReviewAt,
		progress.ProblemID,
		progress.ProblemListID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO progress (perceived_difficulty, status, last_solved_at, next_review_at, problem_id, problem_list_id)
		                VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

		err = r.db.QueryRowContext(ctx, insertQuery,
			progress.PerceivedDifficulty,
			progress.Status,
			progress.LastSolvedAt,
			progress.NextReviewAt,
			progress.ProblemID,
			progress.ProblemListID,
		).Scan(&id)

		if err != nil {
			return 0, fmt.Errorf("failed to insert progress: %w", err)
		}

		return id, nil
	}

	if err != nil {
		return 0, fmt.Errorf("failed to update progress: %w", err)
	}

	return id, nil
}

func (r *progressRepository) GetDueForReview(ctx context.Context) ([]Progress, error) {
	query := `SELECT id, perceived_difficulty, status, last_solved_at, next_review_at, problem_id, problem_list_id
	          FROM progress
	          WHERE status = $1
	            AND next_review_at <= $2
	          ORDER BY next_review_at`

	now := time.Now()
	rows, err := r.db.QueryContext(ctx, query, Active, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get due progress: %w", err)
	}
	defer rows.Close()

	var progressList []Progress
	for rows.Next() {
		var progress Progress
		if err := rows.Scan(
			&progress.ID,
			&progress.PerceivedDifficulty,
			&progress.Status,
			&progress.LastSolvedAt,
			&progress.NextReviewAt,
			&progress.ProblemID,
			&progress.ProblemListID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		progressList = append(progressList, progress)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating progress: %w", err)
	}

	return progressList, nil
}
