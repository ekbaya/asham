package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (h *OrganizationHandler) UpdateNationalTCSecretary(c *gin.Context) {
	var payload struct {
		NSB       string `json:"nsb" binding:"required"`
		Secretary string `json:"secretary" binding:"required"`
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

	err := h.organizationService.UpdateNationalTCSecretary(payload.NSB, payload.Secretary)
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

	// Generate a new UUID for the committee
	newID := uuid.New().String()

	// Add the ID to the committee data
	payload.Committee["id"] = newID

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

func (h *OrganizationHandler) FetchTechnicalCommittees(c *gin.Context) {
	committees, err := h.organizationService.FetchTechnicalCommittees()
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Committee not found")
		return
	}
	c.JSON(http.StatusOK, committees)
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

func (h *OrganizationHandler) SearchTechnicalCommittees(c *gin.Context) {
	// Get the keyword query parameter once
	query := c.Query("keyword")

	// Early return if query is empty to avoid unnecessary service call
	if query == "" {
		c.JSON(http.StatusOK, gin.H{"committees": []any{}})
		return
	}

	// Create search parameters map with only needed fields
	params := map[string]interface{}{
		"name":  query,
		"code":  query,
		"scope": query,
	}

	// Call service to search technical committees
	committees, err := h.organizationService.SearchTechnicalCommittees(params)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"committees": committees,
	})
}

func (h *OrganizationHandler) UpdateCommitteeSecretary(c *gin.Context) {
	var payload struct {
		SecretaryID string `json:"secretary_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	committeeType := c.Param("type")
	committeeID := c.Param("id")

	if err := h.organizationService.UpdateCommitteeSecretary(committeeType, committeeID, payload.SecretaryID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Secretary updated successfully")
}

func (h *OrganizationHandler) UpdateCommitteeChairperson(c *gin.Context) {
	var payload struct {
		ChairpersonID string `json:"chairperson_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	committeeType := c.Param("type")
	committeeID := c.Param("id")

	if err := h.organizationService.UpdateCommitteeChairperson(committeeType, committeeID, payload.ChairpersonID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Chairperson updated successfully")
}

func (h *OrganizationHandler) AddMemberToARSOCouncil(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddMemberToARSOCouncil(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddRegionalEconomicCommunityToJointAdvisoryGroup(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddRegionalEconomicCommunityToJointAdvisoryGroup(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddObserverMemberToJointAdvisoryGroup(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddObserverMemberToJointAdvisoryGroup(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddRegionalRepresentativeToStandardsManagementCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddRegionalRepresentativeToStandardsManagementCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddElectedMemberToStandardsManagementCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddElectedMemberToStandardsManagementCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddObserverToStandardsManagementCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddObserverToStandardsManagementCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddMemberToTechnicalCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddMemberToTechnicalCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddMemberToJointTechnicalCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddMemberToJointTechnicalCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddMemberToSpecializedCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddMemberToSpecializedCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddMemberToWorkingGroup(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddMemberToWorkingGroup(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) AddMemberToTaskForce(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.AddMemberToTaskForce(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member added to Committee")
}

func (h *OrganizationHandler) RemoveMemberFromARSOCouncil(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveMemberFromARSOCouncil(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveRECFromJointAdvisoryGroup(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveRECFromJointAdvisoryGroup(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveObserverFromJointAdvisoryGroup(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveObserverFromJointAdvisoryGroup(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveRegionalElectedMemberFromStandardsManagementCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveRegionalElectedMemberFromStandardsManagementCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveRegionalRepresentativeFromStandardsManagementCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveRegionalRepresentativeFromStandardsManagementCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveMemberFromTechnicalCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveMemberFromTechnicalCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveMemberFromSpecializedCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveMemberFromSpecializedCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) RemoveMemberFromJointTechnicalCommittee(c *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveMemberFromJointTechnicalCommittee(id, payload.MemberID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member removed to Committee")
}

func (h *OrganizationHandler) GetArsoCouncilMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.organizationService.GetArsoCouncilMembers(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *OrganizationHandler) GetJointAdvisoryGroupMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.organizationService.GetJointAdvisoryGroupMembers(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *OrganizationHandler) GetStandardsManagementCommitteeMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.organizationService.GetStandardsManagementCommitteeMembers(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *OrganizationHandler) GetTechnicalCommitteeMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.organizationService.GetTechnicalCommitteeMembers(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *OrganizationHandler) GetSpecializedCommitteeMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.organizationService.GetSpecializedCommitteeMembers(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *OrganizationHandler) GetJointTechnicalCommitteeMembers(c *gin.Context) {
	id := c.Param("id")

	members, err := h.organizationService.GetJointTechnicalCommitteeMembers(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *OrganizationHandler) GetArsoCouncil(c *gin.Context) {
	committees, err := h.organizationService.GetArsoCouncil()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) GetJointAdvisoryGroups(c *gin.Context) {
	committees, err := h.organizationService.GetJointAdvisoryGroups()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) GetStandardsManagementCommittees(c *gin.Context) {
	committees, err := h.organizationService.GetStandardsManagementCommittees()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) GetTechnicalCommittees(c *gin.Context) {
	committees, err := h.organizationService.GetTechnicalCommittees()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) GetSpecializedCommittees(c *gin.Context) {
	committees, err := h.organizationService.GetSpecializedCommittees()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) GetJointTechnicalCommittees(c *gin.Context) {
	committees, err := h.organizationService.GetJointTechnicalCommittees()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) AddMemberStateToTCParticipatingCountries(c *gin.Context) {
	var payload struct {
		TechnicalCommitteeID string `json:"technicalCommitteeId" binding:"required"`
		MemberState          string `json:"memberStateId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.AddMemberStateToTCParticipatingCountries(payload.TechnicalCommitteeID, payload.MemberState)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to add member state")
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member state added successfully")
}

func (h *OrganizationHandler) AddMemberStateToTCObserverCountries(c *gin.Context) {
	var payload struct {
		TechnicalCommitteeID string `json:"technicalCommitteeId" binding:"required"`
		MemberState          string `json:"memberStateId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.AddMemberStateToTCObserverCountries(payload.TechnicalCommitteeID, payload.MemberState)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to add member state")
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member state added successfully")
}

func (h *OrganizationHandler) AddTCToTCEquivalentCommittees(c *gin.Context) {
	var payload struct {
		TechnicalCommitteeID string `json:"technicalCommitteeId" binding:"required"`
		TCToBeAdded          string `json:"newTechnicalCommitteeId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.organizationService.AddTCToTCEquivalentCommittees(payload.TechnicalCommitteeID, payload.TCToBeAdded)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to add technical committee")
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "TC added successfully")
}

func (h *OrganizationHandler) GetTCParticipatingCountries(c *gin.Context) {
	countries, err := h.organizationService.GetTCParticipatingCountries(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"countries": countries})
}

func (h *OrganizationHandler) GetTCObserverCountries(c *gin.Context) {
	countries, err := h.organizationService.GetTCObserverCountries(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"countries": countries})
}

func (h *OrganizationHandler) GetTCEquivalentCommittees(c *gin.Context) {
	committees, err := h.organizationService.GetTCEquivalentCommittees(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"committees": committees})
}

func (h *OrganizationHandler) RemoveMemberStateFromTCParticipatingCountries(c *gin.Context) {
	var payload struct {
		StateID string `json:"state_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveMemberStateFromTCParticipatingCountries(id, payload.StateID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member State removed from tc")
}

func (h *OrganizationHandler) RemoveMemberStateFromTCObserverCountries(c *gin.Context) {
	var payload struct {
		StateID string `json:"state_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveMemberStateFromTCObserverCountries(id, payload.StateID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Member State removed from tc")
}

func (h *OrganizationHandler) RemoveTCFromTCEquivalentCommittees(c *gin.Context) {
	var payload struct {
		TCToBeRemoved string `json:"tc_to_remove"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid input")
		return
	}

	id := c.Param("id")
	if err := h.organizationService.RemoveTCFromTCEquivalentCommittees(id, payload.TCToBeRemoved); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "TC removed from equivalent committees")
}

func (h *OrganizationHandler) GetTCProjects(c *gin.Context) {
	projects, err := h.organizationService.GetTCProjects(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *OrganizationHandler) GetTCWorkingGroups(c *gin.Context) {
	wgs, err := h.organizationService.GetTCWorkingGroups(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"wgs": wgs})
}

func (h *OrganizationHandler) GetCommitteeMeetings(c *gin.Context) {
	meetings, err := h.organizationService.GetCommitteeMeetings(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"meetings": meetings})
}
