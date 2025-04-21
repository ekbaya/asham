package handlers

import (
	"net/http"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AcceptanceHandler struct {
	AcceptanceService services.AcceptanceService
}

func NewAcceptanceHandler(AcceptanceService services.AcceptanceService) *AcceptanceHandler {
	return &AcceptanceHandler{
		AcceptanceService: AcceptanceService,
	}
}

func (h *AcceptanceHandler) CreateNSBResponse(c *gin.Context) {
	var payload models.NSBResponse
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
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	payload.ResponderID = userID.(string)

	err := h.AcceptanceService.CreateNSBResponse(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "NSB response created successfully")
}

func (h *AcceptanceHandler) GetNSBResponse(c *gin.Context) {
	form, err := h.AcceptanceService.GetNSBResponse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, form)
}

func (h *AcceptanceHandler) GetNSBResponsesByProjectID(c *gin.Context) {
	submissions, err := h.AcceptanceService.GetNSBResponsesByProjectID(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, submissions)
}

func (h *AcceptanceHandler) UpdateNSBResponse(c *gin.Context) {
	var payload models.NSBResponse
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

	err := h.AcceptanceService.UpdateNSBResponse(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Record updated successfully")
}

func (h *AcceptanceHandler) DeleteNSBResponse(c *gin.Context) {
	err := h.AcceptanceService.DeleteNSBResponse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Response deleted successfully")
}

func (h *AcceptanceHandler) GetAcceptance(c *gin.Context) {
	acceptance, err := h.AcceptanceService.GetAcceptance(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, acceptance)
}

func (h *AcceptanceHandler) GetAcceptanceByProjectID(c *gin.Context) {
	acceptance, err := h.AcceptanceService.GetAcceptanceByProjectID(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, acceptance)
}

func (h *AcceptanceHandler) GetAcceptances(c *gin.Context) {
	acceptances, err := h.AcceptanceService.GetAcceptances()
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, acceptances)
}

func (h *AcceptanceHandler) UpdateAcceptance(c *gin.Context) {
	var payload models.Acceptance
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

	err := h.AcceptanceService.UpdateAcceptance(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "success")
}

func (h *AcceptanceHandler) GetAcceptanceWithResponses(c *gin.Context) {
	acceptance, err := h.AcceptanceService.GetAcceptanceWithResponses(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, acceptance)
}

func (h *AcceptanceHandler) CountNSBResponsesByType(c *gin.Context) {
	acceptance, err := h.AcceptanceService.CountNSBResponsesByType(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, acceptance)
}

func (h *AcceptanceHandler) CalculateNSBResponseStats(c *gin.Context) {
	err := h.AcceptanceService.CalculateNSBResponseStats(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h *AcceptanceHandler) SetNSBResponseacceptanceApproval(c *gin.Context) {
	var payload models.Acceptance
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

	// Set creator ID from authenticated user
	userID, exists := c.Get("user_id")
	if !exists {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr := userID.(string)
	payload.TCSecretaryID = &userIDStr

	err := h.AcceptanceService.SetAcceptanceApproval(payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Success")
}

func (h *AcceptanceHandler) GetAcceptanceResults(c *gin.Context) {
	results, err := h.AcceptanceService.GetAcceptanceResults(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, results)
}
