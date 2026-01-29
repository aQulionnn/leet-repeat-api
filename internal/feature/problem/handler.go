package problem

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo ProblemRepository
}

func NewHandler(repo ProblemRepository) *Handler {
	return &Handler{repo: repo}
}

// Create godoc
// @Summary Create a new problem
// @Description Create a new problem with question and difficulty
// @Accept json
// @Produce json
// @Param problem body Problem true "Problem object"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /problems [post]
func (h *Handler) Create(c *gin.Context) {
	var problem Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.repo.Add(c.Request.Context(), &problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetAll godoc
// @Summary Get all problems
// @Description Get all problems from database
// @Produce json
// @Success 200 {array} Problem
// @Failure 500 {object} map[string]string
// @Router /problems [get]
func (h *Handler) GetAll(c *gin.Context) {
	problems, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, problems)
}

// GetByID godoc
// @Summary Get problem by ID
// @Description Get a specific problem by ID
// @Produce json
// @Param id path int true "Problem ID"
// @Success 200 {object} Problem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /problems/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	problem, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// Update godoc
// @Summary Update problem
// @Description Update an existing problem
// @Accept json
// @Produce json
// @Param id path int true "Problem ID"
// @Param problem body Problem true "Problem object"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /problems/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var problem Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	problem.ID = id
	rowsAffected, err := h.repo.Update(c.Request.Context(), &problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rows_affected": rowsAffected})
}

// Delete godoc
// @Summary Delete problem
// @Description Delete a problem by ID
// @Produce json
// @Param id path int true "Problem ID"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /problems/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rowsAffected, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rows_affected": rowsAffected})
}
