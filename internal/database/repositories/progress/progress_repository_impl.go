package progress

import (
	"context"
	"fmt"
	"leet-repeat-api/internal/database/models"
	"strings"
)

func (r *progressRepository) BulkUpsert(ctx context.Context, progressList *[]models.Progress) (int, error) {
	if len(*progressList) == 0 {
		return 0, nil
	}

	cols := 8
	placeholders := make([]string, len(*progressList))
	args := make([]interface{}, 0, len(*progressList)*cols)

	for i, p := range *progressList {
		base := i * cols
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8)
		args = append(args, p.PerceivedDifficulty, p.Status, p.LastSolvedAtUtc, p.NextReviewAtUtc, p.ProblemQuestionID, p.ProblemQuestion, p.ProblemDifficulty, p.ProblemListName)
	}

	query := fmt.Sprintf(`
		INSERT INTO progress (perceived_difficulty, status, last_solved_at_utc, next_review_at_utc, problem_question_id, problem_question, problem_difficulty, problem_list_name)
		VALUES %s
		ON CONFLICT (problem_question_id, problem_list_name) 
		DO UPDATE SET
			perceived_difficulty = EXCLUDED.perceived_difficulty,
			status = EXCLUDED.status,
			last_solved_at_utc = EXCLUDED.last_solved_at_utc,
			next_review_at_utc = EXCLUDED.next_review_at_utc
	`, strings.Join(placeholders, ","))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to bulk upsert progressList: %w", err)
	}

	return len(*progressList), nil
}

func (r *progressRepository) GetAll(ctx context.Context) ([]models.Progress, error) {
	query := `
		SELECT 
			id, perceived_difficulty, status, last_solved_at_utc, next_review_at_utc, problem_question_id, problem_question, problem_difficulty, problem_list_name 
		FROM progress
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all progressList: %w", err)
	}

	defer rows.Close()

	var progressList []models.Progress
	for rows.Next() {
		var p models.Progress
		if err := rows.Scan(&p.ID, &p.PerceivedDifficulty, &p.Status, &p.LastSolvedAtUtc, &p.NextReviewAtUtc, &p.ProblemQuestionID, &p.ProblemQuestion, &p.ProblemDifficulty, &p.ProblemListName); err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		progressList = append(progressList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over progress rows: %w", err)
	}

	return progressList, nil
}

func (r *progressRepository) Clear(ctx context.Context) (int, error) {
	query := `DELETE FROM progress`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to clear progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}
