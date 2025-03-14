package handlers

import (
	"errors"
	"net/http"
	"server/db/models"
	"server/db/repository"
	"server/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errUnauthorized = errors.New("unauthorized access")

// BreakdownRequest represents the data needed to create a breakdown
type BreakdownRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type BreakdownHandler struct {
	BaseHandler
	Repo *repository.BreakdownRepository
}

func NewBreakdownHandler(repo *repository.BreakdownRepository) *BreakdownHandler {
	return &BreakdownHandler{
		Repo: repo,
	}
}

// GetBreakdowns retrieves all breakdowns for the authenticated user
func (h *BreakdownHandler) GetBreakdowns(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.HandleError(c, errUnauthorized, http.StatusUnauthorized)
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Find all breakdowns for this user
	var breakdowns []models.Breakdown
	err = h.Repo.Find(c.Request.Context(), bson.M{"user_id": userObjID}, &breakdowns)
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	h.Respond(c, http.StatusOK, breakdowns)
}

// GetBreakdownByID retrieves a specific breakdown by ID
func (h *BreakdownHandler) GetBreakdownByID(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.HandleError(c, errUnauthorized, http.StatusUnauthorized)
		return
	}

	// Parse breakdown ID from URL
	breakdownID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(breakdownID)
	if err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Find the breakdown
	breakdown := &models.Breakdown{}
	err = h.Repo.FindOne(c.Request.Context(), bson.M{
		"_id": objID,
	}, breakdown)
	if err != nil {
		h.HandleError(c, err, http.StatusNotFound)
		return
	}

	// Verify that the breakdown belongs to the authenticated user
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	if breakdown.UserID != userObjID {
		h.HandleError(c, errUnauthorized, http.StatusForbidden)
		return
	}

	h.Respond(c, http.StatusOK, breakdown)
}

// CreateBreakdown creates a new breakdown for the authenticated user
func (h *BreakdownHandler) CreateBreakdown(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.HandleError(c, errUnauthorized, http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request BreakdownRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Convert user ID string to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Create new breakdown
	breakdown := &models.Breakdown{
		Name:        request.Name,
		Description: request.Description,
		UserID:      userObjID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	err = h.Repo.Create(c.Request.Context(), breakdown)
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	h.Respond(c, http.StatusCreated, breakdown)
}

// UpdateBreakdown updates an existing breakdown
func (h *BreakdownHandler) UpdateBreakdown(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.HandleError(c, errUnauthorized, http.StatusUnauthorized)
		return
	}

	// Parse breakdown ID from URL
	breakdownID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(breakdownID)
	if err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Find the breakdown
	existing := &models.Breakdown{}
	err = h.Repo.FindOne(c.Request.Context(), bson.M{"_id": objID}, existing)
	if err != nil {
		h.HandleError(c, err, http.StatusNotFound)
		return
	}

	// Verify that the breakdown belongs to the authenticated user
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	if existing.UserID != userObjID {
		h.HandleError(c, errUnauthorized, http.StatusForbidden)
		return
	}

	// Parse request body
	var request BreakdownRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Update fields
	update := bson.M{
		"$set": bson.M{
			"name":        request.Name,
			"description": request.Description,
			"updated_at":  time.Now(),
		},
	}

	// Save to database
	err = h.Repo.Update(c.Request.Context(), bson.M{"_id": objID}, update)
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// Get the updated breakdown
	updated := &models.Breakdown{}
	err = h.Repo.FindOne(c.Request.Context(), bson.M{"_id": objID}, updated)
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	h.Respond(c, http.StatusOK, updated)
}

// DeleteBreakdown deletes a breakdown
func (h *BreakdownHandler) DeleteBreakdown(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		h.HandleError(c, errUnauthorized, http.StatusUnauthorized)
		return
	}

	// Parse breakdown ID from URL
	breakdownID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(breakdownID)
	if err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Find the breakdown
	existing := &models.Breakdown{}
	err = h.Repo.FindOne(c.Request.Context(), bson.M{"_id": objID}, existing)
	if err != nil {
		h.HandleError(c, err, http.StatusNotFound)
		return
	}

	// Verify that the breakdown belongs to the authenticated user
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	if existing.UserID != userObjID {
		h.HandleError(c, errUnauthorized, http.StatusForbidden)
		return
	}

	// Delete from database
	err = h.Repo.Delete(c.Request.Context(), bson.M{"_id": objID})
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	h.Respond(c, http.StatusOK, gin.H{"message": "Breakdown deleted successfully"})
}
