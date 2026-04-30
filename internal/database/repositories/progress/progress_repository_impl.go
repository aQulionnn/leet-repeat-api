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

	cols := 9
	placeholders := make([]string, len(*progressList))
	args := make([]interface{}, 0, len(*progressList)*cols)

	for i, p := range *progressList {
		base := i * cols
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9)
		args = append(args, p.PerceivedDifficulty, p.Status, p.LastSolvedAtUtc, p.NextReviewAtUtc, p.ProblemQuestionID, p.ProblemQuestion, p.ProblemDifficulty, p.ProblemListName, p.Username)
	}

	query := fmt.Sprintf(`
		INSERT INTO progress (perceived_difficulty, status, last_solved_at_utc, next_review_at_utc, problem_question_id, problem_question, problem_difficulty, problem_list_name, username)
		VALUES %s
		ON CONFLICT (problem_question_id, problem_list_name, username) 
		DO UPDATE SET
			perceived_difficulty = EXCLUDED.perceived_difficulty,
			status = EXCLUDED.status,
			last_solved_at_utc = EXCLUDED.last_solved_at_utc,
			next_review_at_utc = EXCLUDED.next_review_at_utc
		RETURNING id  -- <-- добавь
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to bulk upsert progressList: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		ids = append(ids, id)
	}

	for i, p := range *progressList {
		if len(p.Events) == 0 || i >= len(ids) {
			continue
		}
		progressID := ids[i]

		eventPlaceholders := make([]string, len(p.Events))
		eventArgs := make([]interface{}, 0, len(p.Events)*3)
		for j, e := range p.Events {
			base := j * 3
			eventPlaceholders[j] = fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3)
			eventArgs = append(eventArgs, progressID, e.PerceivedDifficulty, e.SolvedAtUtc)
		}

		eventQuery := fmt.Sprintf(`
			INSERT INTO progress_event (progress_id, perceived_difficulty, solved_at_utc)
			VALUES %s
			ON CONFLICT (progress_id, solved_at_utc) DO NOTHING
		`, strings.Join(eventPlaceholders, ","))

		if _, err := r.db.ExecContext(ctx, eventQuery, eventArgs...); err != nil {
			return 0, fmt.Errorf("failed to insert events: %w", err)
		}
	}

	return len(*progressList), nil
}

func (r *progressRepository) GetAll(ctx context.Context) ([]models.Progress, error) {
	query := `
		SELECT id, perceived_difficulty, status, last_solved_at_utc, next_review_at_utc, 
		       problem_question_id, problem_question, problem_difficulty, problem_list_name, username 
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
		if err := rows.Scan(&p.ID, &p.PerceivedDifficulty, &p.Status, &p.LastSolvedAtUtc, &p.NextReviewAtUtc, &p.ProblemQuestionID, &p.ProblemQuestion, &p.ProblemDifficulty, &p.ProblemListName, &p.Username); err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		progressList = append(progressList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over progress rows: %w", err)
	}

	for i, p := range progressList {
		eventRows, err := r.db.QueryContext(ctx, `
			SELECT id, perceived_difficulty, solved_at_utc 
			FROM progress_event 
			WHERE progress_id = $1
			ORDER BY solved_at_utc DESC
		`, p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get events: %w", err)
		}
		defer eventRows.Close()

		var events []models.ProgressEvent
		for eventRows.Next() {
			var e models.ProgressEvent
			if err := eventRows.Scan(&e.ID, &e.PerceivedDifficulty, &e.SolvedAtUtc); err != nil {
				return nil, fmt.Errorf("failed to scan event: %w", err)
			}
			events = append(events, e)
		}
		progressList[i].Events = events
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
