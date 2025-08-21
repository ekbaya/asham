package handlers

import (
	"net/http"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// BallotingHandler handles HTTP requests related to balloting
type BallotingHandler struct {
	ballotingService *services.BallotingService
}

// NewBallotingHandler creates a new BallotingHandler
func NewBallotingHandler(ballotingService *services.BallotingService) *BallotingHandler {
	return &BallotingHandler{
		ballotingService: ballotingService,
	}
}

// Helper function to extract audit parameters from Gin context
func (h *BallotingHandler) getAuditParams(c *gin.Context) (string, string, string, string, string) {
	userID, exists := c.Get("user_id")
	var userIDStr string
	if exists {
		userIDStr = userID.(string)
	}
	
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	sessionID := c.GetHeader("X-Session-ID")
	requestID := c.GetHeader("X-Request-ID")
	
	return userIDStr, ipAddress, userAgent, sessionID, requestID
}

// CreateBalloting handles the creation of a new balloting session
func (h *BallotingHandler) CreateBalloting(c *gin.Context) {
	var payload models.Balloting
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDStr, ipAddress, userAgent, sessionID, requestID := h.getAuditParams(c)
	err := h.ballotingService.CreateBalloting(&payload, userIDStr, ipAddress, userAgent, sessionID, requestID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "balloting", payload)
}

// GetBallotingByID retrieves a balloting by its ID
func (h *BallotingHandler) GetBallotingByID(c *gin.Context) {
	ballotingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid balloting ID")
		return
	}

	balloting, err := h.ballotingService.FindBallotingByID(ballotingID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "balloting", balloting)
}

// GetAllBallotings retrieves all ballotings
func (h *BallotingHandler) GetAllBallotings(c *gin.Context) {
	ballotings, err := h.ballotingService.FindAll()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "ballotings", ballotings)
}

// GetActiveBallotingsWithVotes retrieves all active ballotings with their votes
func (h *BallotingHandler) GetActiveBallotingsWithVotes(c *gin.Context) {
	ballotings, err := h.ballotingService.FindActiveBallotingWithVotes()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "ballotings", ballotings)
}

// UpdateBalloting updates an existing balloting
func (h *BallotingHandler) UpdateBalloting(c *gin.Context) {
	ballotingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid balloting ID")
		return
	}

	var payload models.Balloting
	if err := c.ShouldBindJSON(&payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure ID matches route parameter
	payload.ID = ballotingID

	userIDStr, ipAddress, userAgent, sessionID, requestID := h.getAuditParams(c)
	err = h.ballotingService.UpdateBalloting(&payload, userIDStr, ipAddress, userAgent, sessionID, requestID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Balloting updated successfully")
}

// DeleteBalloting deletes a balloting
func (h *BallotingHandler) DeleteBalloting(c *gin.Context) {
	ballotingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid balloting ID")
		return
	}

	err = h.ballotingService.DeleteBalloting(ballotingID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Balloting deleted successfully")
}

// GetBallotingsByPeriod retrieves ballotings within a specific time period
func (h *BallotingHandler) GetBallotingsByPeriod(c *gin.Context) {
	var periodQuery struct {
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&periodQuery); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid date range format")
		return
	}

	startDate, err := time.Parse("2006-01-02", periodQuery.StartDate)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid start date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", periodQuery.EndDate)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid end date format")
		return
	}

	// Add one day to end date to include the entire end date
	endDate = endDate.Add(24 * time.Hour)

	ballotings, err := h.ballotingService.FindBallotingByPeriod(startDate, endDate)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "ballotings", ballotings)
}

func (h *BallotingHandler) RecommendFDARS(c *gin.Context) {
	var payload struct {
		Project     string `json:"project_id" binding:"required"`
		Recommended bool   `json:"recommended" `
	}

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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	err := h.ballotingService.RecommendFDARS(userIDStr, payload.Project, payload.Recommended)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Recommendation updated successfully")
}

func (h *BallotingHandler) VerifyFDARSRecommendation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)
	err := h.ballotingService.VerifyFDARSRecommendation(userIDStr, c.Param("project_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Recommendation verified updated successfully")
}
