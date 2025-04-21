package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MeetingHandler struct {
	meetingService services.MeetingService
}

func NewMeetingHandler(meetingService services.MeetingService) *MeetingHandler {
	return &MeetingHandler{
		meetingService: meetingService,
	}
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	var payload models.Meeting
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert string to UUID
	createdByID, err := uuid.Parse(userID.(string))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Set created by
	payload.CreatedByID = createdByID.String()
	payload.CreatedAt = time.Now()

	// Validate meeting date
	if payload.Date.Before(time.Now()) {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting date cannot be in the past")
		return
	}

	// Check if agenda is prepared at least 3 weeks before meeting
	threeWeeks := 21 * 24 * time.Hour
	timeUntilMeeting := payload.Date.Sub(time.Now())
	if timeUntilMeeting < threeWeeks && payload.Agenda == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Agenda must be provided at least 3 weeks before meeting date")
		return
	}

	// Set initial status
	payload.Status = models.MeetingStatusPlanned

	err = h.meetingService.CreateMeeting(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Meeting created successfully")
}

func (h *MeetingHandler) GetMeetingByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	meeting, err := h.meetingService.GetMeetingByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if meeting == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Meeting not found")
		return
	}

	utilities.Show(c, http.StatusOK, "meeting", meeting)
}

func (h *MeetingHandler) GetAllMeetings(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid page size")
		return
	}

	meetings, err := h.meetingService.GetAllMeetings(page, pageSize)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meetings": meetings,
		"total":    len(*meetings),
		"limit":    pageSize,
		"page":     page,
	})

}

func (h *MeetingHandler) GetMeetingsByCommittee(c *gin.Context) {
	committeeID := c.Param("committee_id")
	if committeeID == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Committee ID is required")
		return
	}

	meetings, err := h.meetingService.GetMeetingsByCommittee(committeeID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "meetings", meetings)
}

func (h *MeetingHandler) GetMeetingsByProject(c *gin.Context) {
	projectID := c.Param("project_id")
	if projectID == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Project ID is required")
		return
	}

	meetings, err := h.meetingService.GetMeetingsByProject(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "meetings", meetings)
}

func (h *MeetingHandler) GetUpcomingMeetings(c *gin.Context) {
	meetings, err := h.meetingService.GetUpcomingMeetings()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "meetings", meetings)
}

func (h *MeetingHandler) UpdateMeeting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	// First get the existing meeting
	existingMeeting, err := h.meetingService.GetMeetingByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingMeeting == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Meeting not found")
		return
	}

	// Then parse the update payload
	var payload models.Meeting
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

	// Maintain original values that shouldn't be changed
	payload.ID = existingMeeting.ID
	payload.CreatedBy = existingMeeting.CreatedBy
	payload.CreatedAt = existingMeeting.CreatedAt

	// Validate meeting date change
	if payload.Date.Before(time.Now()) {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting date cannot be in the past")
		return
	}

	err = h.meetingService.UpdateMeeting(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Meeting updated successfully")
}

func (h *MeetingHandler) DeleteMeeting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	// Check if the meeting exists
	meeting, err := h.meetingService.GetMeetingByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if meeting == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Meeting not found")
		return
	}

	// Don't allow deletion of meetings within 2 weeks
	twoWeeks := 14 * 24 * time.Hour
	if time.Until(meeting.Date) < twoWeeks {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meetings cannot be deleted within two weeks of scheduled date")
		return
	}

	err = h.meetingService.DeleteMeeting(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Meeting deleted successfully")
}

func (h *MeetingHandler) UpdateMeetingStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	var payload struct {
		Status models.MeetingStatus `json:"status" binding:"required"`
		Reason string               `json:"reason"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.meetingService.UpdateMeetingStatus(id, payload.Status, payload.Reason)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Meeting status updated successfully")
}

func (h *MeetingHandler) AddAttendeeToMeeting(c *gin.Context) {
	meetingID := c.Param("meeting_id")
	if meetingID == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	var payload struct {
		MemberID string `json:"member_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.meetingService.AddAttendeeToMeeting(meetingID, payload.MemberID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Attendee added successfully")
}

func (h *MeetingHandler) RemoveAttendeeFromMeeting(c *gin.Context) {
	meetingID := c.Param("meeting_id")
	memberID := c.Param("member_id")

	if meetingID == "" || memberID == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID and Member ID are required")
		return
	}

	err := h.meetingService.RemoveAttendeeFromMeeting(meetingID, memberID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Attendee removed successfully")
}

func (h *MeetingHandler) AddRelatedDocumentToMeeting(c *gin.Context) {
	meetingID := c.Param("meeting_id")
	if meetingID == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	var payload struct {
		DocumentID string `json:"document_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.meetingService.AddRelatedDocumentToMeeting(meetingID, payload.DocumentID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Document added successfully")
}

func (h *MeetingHandler) CheckQuorum(c *gin.Context) {
	meetingID := c.Param("meeting_id")
	if meetingID == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Meeting ID is required")
		return
	}

	hasQuorum, err := h.meetingService.CheckQuorum(meetingID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "has_quorum", hasQuorum)
}
