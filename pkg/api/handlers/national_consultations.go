package handlers

import (
	"net/http"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// NationalConsultationHandler struct
type NationalConsultationHandler struct {
	NationalConsultationService *services.NationalConsultationService
}

// NewNationalConsultationHandler initializes a new NationalConsultationHandler
func NewNationalConsultationHandler(NationalConsultationService *services.NationalConsultationService) *NationalConsultationHandler {
	return &NationalConsultationHandler{
		NationalConsultationService: NationalConsultationService,
	}
}

func (h *NationalConsultationHandler) CreateNationalConsultation(c *gin.Context) {
	var payload models.NationalConsultation
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

	payload.NationalSecretaryID = userID.(string)

	err := h.NationalConsultationService.Create(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "NationalConsultation added successfully")
}

func (h *NationalConsultationHandler) GetNationalConsultationByID(c *gin.Context) {
	NationalConsultationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid NationalConsultation ID")
		return
	}

	NationalConsultation, err := h.NationalConsultationService.GetByID(NationalConsultationID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "NationalConsultation", NationalConsultation)
}

func (h *NationalConsultationHandler) GetAllNationalConsultations(c *gin.Context) {
	NationalConsultations, err := h.NationalConsultationService.GetAll()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "NationalConsultations", NationalConsultations)
}

func (h *NationalConsultationHandler) UpdateNationalConsultation(c *gin.Context) {
	NationalConsultationID, err := uuid.Parse(c.Param("NationalConsultation_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid NationalConsultation ID")
		return
	}

	var payload models.NationalConsultation
	if err := c.ShouldBindJSON(&payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	payload.ID = NationalConsultationID // Ensure ID matches request parameter
	err = h.NationalConsultationService.Update(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "NationalConsultation updated successfully")
}

func (h *NationalConsultationHandler) DeleteNationalConsultation(c *gin.Context) {
	NationalConsultationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid NationalConsultation ID")
		return
	}

	err = h.NationalConsultationService.Delete(NationalConsultationID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "NationalConsultation deleted successfully")
}

func (h *NationalConsultationHandler) GetNationalConsultationsByProjectID(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	NationalConsultations, err := h.NationalConsultationService.GetByProjectID(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "NationalConsultations", NationalConsultations)
}
