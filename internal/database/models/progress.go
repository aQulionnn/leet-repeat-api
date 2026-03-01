package models

import (
	"leet-repeat-api/internal/database/enums/difficulty"
	"leet-repeat-api/internal/database/enums/perceived_difficulty"
	"leet-repeat-api/internal/database/enums/status"
	"time"
)

type Progress struct {
	ID                  int                                      `json:"id"`
	PerceivedDifficulty perceived_difficulty.PerceivedDifficulty `json:"perceived_difficulty"`
	Status              status.Status                            `json:"status"`
	LastSolvedAtUtc     *time.Time                               `json:"last_solved_at_utc"`
	NextReviewAtUtc     *time.Time                               `json:"next_review_at_utc"`
	ProblemQuestionID   int                                      `json:"problem_question_id"`
	ProblemQuestion     string                                   `json:"problem_question"`
	ProblemDifficulty   difficulty.Difficulty                    `json:"problem_difficulty"`
	ProblemListName     string                                   `json:"problem_list_name"`
}
