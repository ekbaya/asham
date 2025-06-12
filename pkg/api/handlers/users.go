package handlers

import (
	"net/http"
	"strconv"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UsersHandler struct {
	userService services.MemberService
}

// Constructor now takes the service as a parameter
func NewUsersHandler(userService services.MemberService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
	}
}

func (h *UsersHandler) RegisterMember(c *gin.Context) {
	var payload models.Member
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Convert validation errors into human-readable messages
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}

		// For non-validation errors
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.userService.CreateMember(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "User registered successfully")
}

func (h *UsersHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind and validate the request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Convert validation errors into human-readable messages
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}

		// For non-validation errors
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Authenticate user and generate tokens
	token, refreshToken, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		// Handle authentication errors (e.g., invalid credentials)
		utilities.ShowMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Return the tokens in the response
	c.JSON(http.StatusOK, gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
		"expires_in":    86400,
	})
}

func (h *UsersHandler) GenerateRefreshToken(c *gin.Context) {
	// Retrieve user_id from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Retrieve user details
	user, err := h.userService.Account(userID.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found or unauthorized"})
		return
	}

	// Generate a new access token
	accessToken, err := models.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate a new refresh token
	refreshToken, err := models.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Return the new tokens in the response
	c.JSON(http.StatusOK, gin.H{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"expires_in":         86400,  // Access token expiration in seconds (24 hours)
		"refresh_expires_in": 604800, // Refresh token expiration in seconds (7 days)
	})
}

func (h *UsersHandler) Account(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.AccountWithResponsibilities(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	utilities.Show(c, http.StatusOK, "account", user)
}

func (h *UsersHandler) GetUserDetails(c *gin.Context) {
	user, err := h.userService.AccountWithResponsibilities(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UsersHandler) UpdateUser(c *gin.Context) {
	var payload models.Member
	c.ShouldBindJSON(&payload)

	err := h.userService.UpdateMember(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Member updated successfully")
}

func (h *UsersHandler) GetAllUsers(c *gin.Context) {

	limit := 10
	offset := 0

	if limitStr := c.Query("pageSize"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("page"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	members, total, err := h.userService.GetAllMembers(limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": members,
		"total": total,
		"limit": limit,
		"page":  offset,
	})
}

func (h *UsersHandler) DeleteMember(c *gin.Context) {
	err := h.userService.DeleteMember(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusUnauthorized, err.Error())
		return
	}
	utilities.ShowMessage(c, http.StatusNoContent, "Member deleted successfully")
}

func (h *UsersHandler) LogoutAll(c *gin.Context) {
	err := h.userService.LogoutAll(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "User logged out")
}
