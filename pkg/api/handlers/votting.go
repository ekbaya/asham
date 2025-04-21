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

// VoteHandler handles HTTP requests related to votes
type VoteHandler struct {
	ballotingService *services.BallotingService
}

// NewVoteHandler creates a new VoteHandler
func NewVoteHandler(ballotingService *services.BallotingService) *VoteHandler {
	return &VoteHandler{
		ballotingService: ballotingService,
	}
}

// CreateVote handles the creation of a new vote
func (h *VoteHandler) CreateVote(c *gin.Context) {
	var payload models.Vote
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

	// Set ID if not provided
	if payload.ID == uuid.Nil {
		payload.ID = uuid.New()
	}

	// Set voter ID from authenticated user
	userID, exists := c.Get("user_id")
	if exists {
		payload.MemberID = userID.(string)
	}

	// Set creation time
	payload.CreatedAt = time.Now()

	err := h.ballotingService.CreateVote(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "vote", payload)
}

// GetVoteByID retrieves a vote by its ID
func (h *VoteHandler) GetVoteByID(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid vote ID")
		return
	}

	vote, err := h.ballotingService.FindVoteByID(voteID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "vote", vote)
}

// GetVotesByBallotingID retrieves all votes for a specific balloting
func (h *VoteHandler) GetVotesByBallotingID(c *gin.Context) {
	ballotingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid balloting ID")
		return
	}

	votes, err := h.ballotingService.FindVotesByBallotingID(ballotingID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "votes", votes)
}

// GetVotesByProjectID retrieves all votes for a specific project
func (h *VoteHandler) GetVotesByProjectID(c *gin.Context) {
	projectID := c.Param("id")

	votes, err := h.ballotingService.FindByProjectID(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "votes", votes)
}

// GetVotesByMemberID retrieves all votes cast by a specific member
func (h *VoteHandler) GetVotesByMemberID(c *gin.Context) {
	memberID := c.Param("member_id")

	votes, err := h.ballotingService.FindVotesByMemberID(memberID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "votes", votes)
}

// UpdateVote updates an existing vote
func (h *VoteHandler) UpdateVote(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid vote ID")
		return
	}

	// First retrieve the existing vote to preserve fields not in payload
	existingVote, err := h.ballotingService.FindVoteByID(voteID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	var payload models.Vote
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
	payload.ID = voteID

	// Preserve fields that shouldn't be updated
	payload.MemberID = existingVote.MemberID
	payload.ProjectID = existingVote.ProjectID
	payload.BallotingID = existingVote.BallotingID
	payload.CreatedAt = existingVote.CreatedAt

	err = h.ballotingService.UpdateVote(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Vote updated successfully")
}

// DeleteVote deletes a vote
func (h *VoteHandler) DeleteVote(c *gin.Context) {
	voteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid vote ID")
		return
	}

	err = h.ballotingService.DeleteVote(voteID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Vote deleted successfully")
}

// GetAllVotesWithAssociations retrieves all votes with their related data
func (h *VoteHandler) GetAllVotesWithAssociations(c *gin.Context) {
	votes, err := h.ballotingService.FindVotesWithAssociations()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "votes", votes)
}

// CountVotesByBalloting counts votes for a specific balloting
func (h *VoteHandler) CountVotesByBalloting(c *gin.Context) {
	ballotingID, err := uuid.Parse(c.Param("balloting_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid balloting ID")
		return
	}

	count, err := h.ballotingService.CountVotesByBalloting(ballotingID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "vote_count", count)
}

// CheckProjectAcceptanceCriteria checks if a project meets the acceptance criteria
func (h *VoteHandler) CheckProjectAcceptanceCriteria(c *gin.Context) {
	projectID := c.Param("id")

	result, err := h.ballotingService.CheckAcceptanceCriteria(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "acceptance_criteria", result)
}
