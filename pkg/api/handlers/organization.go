package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrganizationHandler struct {
	organizationService services.OrganizationService
}

func NewOrganizationHandler(organizationService services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: organizationService,
	}
}

func (h *OrganizationHandler) CreateMemberState(c *gin.Context) {
	var payload models.MemberState
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Convert validation errors into human-readable messages
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}

		// For non-validation errors
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CreateMemberState(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Member state added successfully")
}

func (h *OrganizationHandler) FetchMemberStates(c *gin.Context) {
	states, err := h.organizationService.FetchMemberStates()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, states)
}

func (h *OrganizationHandler) CreateNSB(c *gin.Context) {
	var payload models.NationalStandardBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Convert validation errors into human-readable messages
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}

		// For non-validation errors
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CreateNSB(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "NSB registered successfully")
}

func (h *OrganizationHandler) FetchNSBs(c *gin.Context) {
	nsbs, err := h.organizationService.FetchNSBs()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nsbs)
}

func (h *OrganizationHandler) CreateCommittee(c *gin.Context) {
	var payload struct {
		Type      models.CommitteeType `json:"type" binding:"required"`
		Committee map[string]any       `json:"committee" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate committee type using the enum
	if !models.ValidateCommitteeType(string(payload.Type)) {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid committee type")
		return
	}

	// Map to the correct struct based on Type
	var committee any
	switch models.CommitteeType(payload.Type) {
	case models.ARSO_Council:
		committee = &models.ARSOCouncil{}
	case models.Joint_Advisory_Group:
		committee = &models.JointAdvisoryGroup{}
	case models.Standards_Management_Committee:
		committee = &models.StandardsManagementCommittee{}
	case models.Technical_Committee:
		committee = &models.TechnicalCommittee{}
	case models.Specialized_Committee:
		committee = &models.SpecializedCommittee{}
	case models.Joint_Technical_Committee:
		committee = &models.JointTechnicalCommittee{}
	}

	// Convert map to struct
	data, err := json.Marshal(payload.Committee)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to process committee data")
		return
	}

	if err := json.Unmarshal(data, committee); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid committee data format")
		return
	}

	// Validate committee struct
	validate := validator.New()
	if err := validate.Struct(committee); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create committee
	err = h.organizationService.CreateCommittee(committee)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Committee registered successfully")
}

func (h *OrganizationHandler) GetCommitteeByID(c *gin.Context) {
	id := c.Param("id")
	committeeType := c.Param("type")

	// Map to the correct struct based on Type
	var model any
	switch models.CommitteeType(committeeType) {
	case models.ARSO_Council:
		model = &models.ARSOCouncil{}
	case models.Joint_Advisory_Group:
		model = &models.JointAdvisoryGroup{}
	case models.Standards_Management_Committee:
		model = &models.StandardsManagementCommittee{}
	case models.Technical_Committee:
		model = &models.TechnicalCommittee{}
	case models.Specialized_Committee:
		model = &models.SpecializedCommittee{}
	case models.Joint_Technical_Committee:
		model = &models.JointTechnicalCommittee{}
	}

	committee, err := h.organizationService.GetCommitteeByID(id, model)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, committee)
}

func (h *OrganizationHandler) UpdateCommittee(c *gin.Context) {
	var payload struct {
		Type      models.CommitteeType `json:"type" binding:"required"`
		Committee map[string]any       `json:"committee" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate committee type using the enum
	if !models.ValidateCommitteeType(string(payload.Type)) {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid committee type")
		return
	}

	// Map to the correct struct based on Type
	var committee any
	switch models.CommitteeType(payload.Type) {
	case models.ARSO_Council:
		committee = &models.ARSOCouncil{}
	case models.Joint_Advisory_Group:
		committee = &models.JointAdvisoryGroup{}
	case models.Standards_Management_Committee:
		committee = &models.StandardsManagementCommittee{}
	case models.Technical_Committee:
		committee = &models.TechnicalCommittee{}
	case models.Specialized_Committee:
		committee = &models.SpecializedCommittee{}
	case models.Joint_Technical_Committee:
		committee = &models.JointTechnicalCommittee{}
	}

	// Convert map to struct
	data, err := json.Marshal(payload.Committee)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to process committee data")
		return
	}

	if err := json.Unmarshal(data, committee); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid committee data format")
		return
	}

	// Validate committee struct
	validate := validator.New()
	if err := validate.Struct(committee); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update committee
	err = h.organizationService.UpdateCommittee(committee)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Committee updated successfully")
}

func (h *OrganizationHandler) DeleteCommittee(c *gin.Context) {
	id := c.Param("id")
	committeeType := c.Param("type")

	err := h.organizationService.DeleteCommittee(committeeType, id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, "Committee deleted succussfuly")
}

func (h *OrganizationHandler) AddWorkingGroupToTechnicalCommittee(c *gin.Context) {
	var payload struct {
		TechnicalCommitteeID string              `json:"technicalCommitteeId" binding:"required"`
		WorkingGroup         models.WorkingGroup `json:"workingGroup" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	tc, err := h.organizationService.GetCommitteeByID(payload.TechnicalCommitteeID, &models.TechnicalCommittee{})
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Technical Committee not found")
		return
	}

	err = h.organizationService.AddWorkingGroupToTechnicalCommittee(tc.(*models.TechnicalCommittee), &payload.WorkingGroup)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Working Group added successfully")
}

func (h *OrganizationHandler) CompleteWorkingGroup(c *gin.Context) {
	var wg models.WorkingGroup
	if err := c.ShouldBindJSON(&wg); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CompleteWorkingGroup(&wg)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Working Group completed successfully")
}

func (h *OrganizationHandler) CreateWorkingGroup(c *gin.Context) {
	var wg models.WorkingGroup
	if err := c.ShouldBindJSON(&wg); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CreateWorkingGroup(&wg)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Working Group created successfully")
}

func (h *OrganizationHandler) GetWorkingGroupByID(c *gin.Context) {
	id := c.Param("id")

	wg, err := h.organizationService.GetWorkingGroupByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Working Group not found")
		return
	}

	c.JSON(http.StatusOK, wg)
}

func (h *OrganizationHandler) GetCommitteeWorkingGroups(c *gin.Context) {
	id := c.Param("id")

	wg, err := h.organizationService.GetCommitteeWorkingGroups(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Working Group not found")
		return
	}

	c.JSON(http.StatusOK, wg)
}

func (h *OrganizationHandler) CreateTaskForce(c *gin.Context) {
	var tf models.TaskForce
	if err := c.ShouldBindJSON(&tf); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CreateTaskForce(&tf)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Task Force created successfully")
}

func (h *OrganizationHandler) GetTaskForceByID(c *gin.Context) {
	id := c.Param("id")

	tf, err := h.organizationService.GetTaskForceByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Task Force not found")
		return
	}

	c.JSON(http.StatusOK, tf)
}

func (h *OrganizationHandler) CreateSubCommittee(c *gin.Context) {
	var sc models.SubCommittee
	if err := c.ShouldBindJSON(&sc); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CreateSubCommittee(&sc)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Sub-Committee created successfully")
}

func (h *OrganizationHandler) AddMemberToSubCommittee(c *gin.Context) {
	var payload struct {
		SubCommitteeID string        `json:"subCommitteeId" binding:"required"`
		Member         models.Member `json:"member" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	sc, err := h.organizationService.GetCommitteeByID(payload.SubCommitteeID, &models.SubCommittee{})
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Sub-Committee not found")
		return
	}

	err = h.organizationService.AddMemberToSubCommittee(sc.(*models.SubCommittee), &payload.Member)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Sub-Committee successfully")
}

func (h *OrganizationHandler) CreateSpecializedCommittee(c *gin.Context) {
	var sc models.SpecializedCommittee
	if err := c.ShouldBindJSON(&sc); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.CreateSpecializedCommittee(&sc)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Specialized Committee created successfully")
}

func (h *OrganizationHandler) GetSpecializedCommitteeByType(c *gin.Context) {
	typeParam := c.Param("type")

	sc, err := h.organizationService.GetSpecializedCommitteeByType(typeParam)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Specialized Committee not found")
		return
	}

	c.JSON(http.StatusOK, sc)
}
