package progress

import "time"

type Progress struct {
	ID                  int                 `json:"id"`
	PerceivedDifficulty PerceivedDifficulty `json:"perceived_difficulty"`
	Status				Status				`json:"status"`
	LastSolvedAt        *time.Time          `json:"last_solved_at"`
	NextReviewAt        *time.Time          `json:"next_review_at"`
	ProblemID           int                 `json:"problem_id"`
	ProblemListID       int                 `json:"problem_list_id"`
}

type PerceivedDifficulty int

const (
	VeryEasy PerceivedDifficulty = iota
	Easy
	Medium
	Hard
	VeryHard
	ExtremelyHard
)

type Status int

const (
	Active Status = iota
	Mastered
	Paused
)
