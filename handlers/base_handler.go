package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseHandler provides common response methods for all handlers.
type BaseHandler struct{}

// Respond sends a JSON response with a status code.
func (h *BaseHandler) Respond(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// HandleError sends a standardized error response.
func (h *BaseHandler) HandleError(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}

// Create handles the creation of resources (a generic method for creating documents).
func (h *BaseHandler) Create(c *gin.Context, repo interface{}, document interface{}) {
	// Bind JSON payload into the document (i.e., breakdown model, user model, etc.)
	if err := c.ShouldBindJSON(document); err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// // Assuming the repository has a `Create` method (this could be customized further)
	// if err := repo.Create(c.Request.Context(), document); err != nil {
	// 	h.HandleError(c, err, http.StatusInternalServerError)
	// 	return
	// }

	h.Respond(c, http.StatusCreated, gin.H{"message": "Resource created successfully"})
}
