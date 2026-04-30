package handlers

import (
	"leet-repeat-api/internal/database/enums/difficulty"
	"leet-repeat-api/internal/database/enums/perceived_difficulty"
	"leet-repeat-api/internal/database/enums/status"
	"leet-repeat-api/internal/database/models"
	"leet-repeat-api/internal/database/repositories/progress"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	repo progress.ProgressRepository
}

func NewProgressHandler(repo progress.ProgressRepository) *ProgressHandler {
	return &ProgressHandler{repo: repo}
}

var (
	cache      = &sync.Map{}
	cacheTimer *time.Timer
	timerMu    sync.Mutex
)

// @Summary Bulk upsert progress
// @Description Insert or update multiple progress records. Conflicts on (problemQuestionId, problemListName) are resolved by updating the existing record.
// @Tags progress
// @Accept json
// @Produce json
// @Param progressList body []bulkUpsertRequest true "List of progress records"
// @Success 200 {object} bulkUpsertResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/progress/bulk-upsert [post]
func (h *ProgressHandler) BulkUpsert(c *gin.Context) {
	var request []bulkUpsertRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Invalid request body"})
		return
	}

	progressList := mapToProgressList(request)

	count, err := h.repo.BulkUpsert(c.Request.Context(), &progressList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	progress, _ := h.repo.GetAll(c.Request.Context())
	cache.Store("progress", progress)

	timerMu.Lock()
	if cacheTimer != nil {
		cacheTimer.Stop()
	}
	cacheTimer = time.AfterFunc(15 * time.Minute, func() {
		cache.Delete("progress")
	})
	timerMu.Unlock()

	c.JSON(http.StatusOK, bulkUpsertResponse{
		Message: "progress list upserted successfully",
		Count:   count,
	})
}

// @Summary Get all progress
// @Description Returns all progress records
// @Tags progress
// @Produce json
// @Success 200 {array} models.Progress
// @Failure 500 {object} errorResponse
// @Router /api/progress [get]
func (h *ProgressHandler) GetAll(c *gin.Context) {
	if cached, ok := cache.Load("progress"); ok {
		c.JSON(http.StatusOK, cached)
		return
	}

	progressList, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progressList)
}

// @Summary Clear all progress
// @Description Deletes all progress records
// @Tags progress
// @Produce json
// @Success 200 {object} clearResponse
// @Failure 500 {object} errorResponse
// @Router /api/progress/clear [delete]
func (h *ProgressHandler) Clear(c *gin.Context) {
	count, err := h.repo.Clear(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, ok := cache.Load("progress"); ok {
		cache.Clear()
	}

	c.JSON(http.StatusOK, gin.H{"message": "progress cleared successfully", "count": count})
}

func mapToProgressList(requests []bulkUpsertRequest) []models.Progress {
	var progressList []models.Progress
	for _, item := range requests {
		events := make([]models.ProgressEvent, len(item.Events))
		for i, e := range item.Events {
			events[i] = models.ProgressEvent{
				PerceivedDifficulty: e.PerceivedDifficulty,
				SolvedAtUtc:         e.SolvedAtUtc,
			}
		}

		progressList = append(progressList, models.Progress{
			PerceivedDifficulty: perceived_difficulty.PerceivedDifficulty(item.PerceivedDifficulty),
			Status:              status.Status(item.Status),
			LastSolvedAtUtc:     item.LastSolvedAtUtc,
			NextReviewAtUtc:     item.NextReviewAtUtc,
			ProblemQuestionID:   item.ProblemQuestionID,
			ProblemQuestion:     item.ProblemQuestion,
			ProblemDifficulty:   difficulty.Difficulty(item.ProblemDifficulty),
			ProblemListName:     item.ProblemListName,
			Username:            item.Username,
			Events:              events,
		})
	}
	return progressList
}

type bulkUpsertRequest struct {
	PerceivedDifficulty int                      `json:"perceivedDifficulty" example:"0"`
	Username            string                   `json:"username"            example:"john_doe"`
	Status              int                      `json:"status"              example:"0"`
	LastSolvedAtUtc     *time.Time               `json:"lastSolvedAtUtc"     example:"2025-01-01T10:00:00Z"`
	NextReviewAtUtc     *time.Time               `json:"nextReviewAtUtc"     example:"2025-01-02T10:00:00Z"`
	ProblemQuestionID   int                      `json:"problemQuestionId"   example:"1"`
	ProblemQuestion     string                   `json:"problemQuestion"     example:"Two Sum"`
	ProblemDifficulty   int                      `json:"problemDifficulty"   example:"0"`
	ProblemListName     string                   `json:"problemListName"     example:"Arrays"`
	Events              []bulkUpsertRequestEvent `json:"events"`
}

type bulkUpsertRequestEvent struct {
	PerceivedDifficulty int        `json:"perceivedDifficulty" example:"0"`
	SolvedAtUtc         *time.Time `json:"solvedAtUtc"         example:"2025-01-01T10:00:00Z"`
}

type bulkUpsertResponse struct {
	Message string `json:"message" example:"progress list upserted successfully"`
	Count   int    `json:"count"   example:"5"`
}

type clearResponse struct {
	Message string `json:"message" example:"progress cleared successfully"`
	Count   int    `json:"count"   example:"3"`
}

type errorResponse struct {
	Error string `json:"error" example:"Invalid request body"`
}
