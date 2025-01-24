package handlers

import (
	"net/http"
	"server/db/models"
	"server/db/repository"

	"github.com/gin-gonic/gin"
)

type BreakdownHandler struct {
	BaseHandler
	Repo *repository.BreakdownRepository
}

func NewBreakdownHandler(repo *repository.BreakdownRepository) *BreakdownHandler {
	return &BreakdownHandler{
		Repo: repo,
	}
}

func (h *BreakdownHandler) GetBreakdowns(c *gin.Context) {
	var breakdowns []interface{}
	// Use the base repository's GetAll method
	err := h.Repo.FindAll(c.Request.Context(), &breakdowns)
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}
	h.Respond(c, http.StatusOK, breakdowns)
}

// CreateBreakdownsHandler creates a new breakdown.
func (h *BreakdownHandler) CreateBreakdowns(c *gin.Context) {
	breakdown := &models.Breakdown{}
	h.Create(c, h.Repo, breakdown)
}
