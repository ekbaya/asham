package handlers

import (
	"net/http"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StandardHandler struct {
	standardService *services.StandardService
}

func NewStandardHandler(standardService *services.StandardService) *StandardHandler {
	return &StandardHandler{
		standardService: standardService,
	}
}

// Create a new standard
func (h *StandardHandler) CreateStandard(c *gin.Context) {
	var payload models.Standard
	if err := c.ShouldBindJSON(&payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formatted := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formatted)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	payload.CreatedAt = time.Now()
	payload.UpdatedAt = payload.CreatedAt
	payload.Version = 1
	payload.UpdatedByID = userIDStr

	if err := h.standardService.CreateStandard(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "success", payload)
}

// Save standard (Auto-save / webhook-style)
func (h *StandardHandler) SaveStandard(c *gin.Context) {
	id := c.Param("id")
	var payload models.Standard

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	standard, err := h.standardService.GetStandardByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Standard not found")
		return
	}

	standard.Content = payload.Content
	standard.UpdatedBy = payload.UpdatedBy
	standard.UpdatedAt = time.Now()

	if err := h.standardService.SaveStandard(standard, userIDStr); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Standard saved successfully")
}

// Get a standard by ID
func (h *StandardHandler) GetStandard(c *gin.Context) {
	id := c.Param("id")

	standard, err := h.standardService.GetStandardByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Standard not found")
		return
	}

	utilities.Show(c, http.StatusOK, "standard", standard)
}

// Get version history of a standard
func (h *StandardHandler) GetStandardVersions(c *gin.Context) {
	id := c.Param("id")

	versions, err := h.standardService.GetStandardVersions(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "versions", versions)
}

// Restore a specific version
func (h *StandardHandler) RestoreVersion(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Version int `json:"version"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.standardService.RestoreStandardVersion(id, payload.Version); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Standard version restored")
}

func (h *StandardHandler) DiffVersions(c *gin.Context) {
	id := c.Param("id")
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Both from and to versions are required")
		return
	}

	standardID := id

	fromVersion, err := h.standardService.GetVersion(standardID, from)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "From version not found")
		return
	}

	toVersion, err := h.standardService.GetVersion(standardID, to)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "To version not found")
		return
	}

	diff, err := utilities.DiffJSON(fromVersion.Content, toVersion.Content)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to diff versions")
		return
	}

	utilities.Show(c, http.StatusOK, "diff", diff)
}

func (h *StandardHandler) GetAuditLogs(c *gin.Context) {
	standardID := c.Param("id")

	logs, err := h.standardService.GetAuditLogs(standardID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to fetch audit logs")
		return
	}

	utilities.Show(c, http.StatusOK, "logs", logs)
}
