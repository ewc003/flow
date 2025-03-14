package handlers

import (
	"net/http"
	"server/db/models"
	"server/db/repository"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler
	UserRepo *repository.UserRepository
}

// LoginRequest represents the login form data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration form data
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo: repo,
	}
}

// Login handles user authentication and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Validate credentials
	user, err := h.UserRepo.ValidateCredentials(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		h.HandleError(c, err, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// Return the token
	h.Respond(c, http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Create new user
	user := &models.User{
		Username:  request.Username,
		Email:     request.Email,
		Password:  request.Password, // Will be hashed in repository
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	err := h.UserRepo.CreateUser(c.Request.Context(), user)
	if err != nil {
		h.HandleError(c, err, http.StatusBadRequest)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		h.HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// Return the token
	h.Respond(c, http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetProfile retrieves the authenticated user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		h.HandleError(c, errUnauthorized, http.StatusUnauthorized)
		return
	}

	// Find user by ID
	user, err := h.UserRepo.FindUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.HandleError(c, err, http.StatusNotFound)
		return
	}

	// Return user profile (exclude password)
	h.Respond(c, http.StatusOK, gin.H{
		"id":        user.ID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}
