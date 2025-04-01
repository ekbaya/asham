package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProposalHandler struct {
	proposalService *services.ProposalService
	documentService *services.DocumentService // For handling referenced standards
	projectService  *services.ProjectService  // For project validation
}

func NewProposalHandler(
	proposalService *services.ProposalService,
	documentService *services.DocumentService,
	projectService *services.ProjectService,
) *ProposalHandler {
	return &ProposalHandler{
		proposalService: proposalService,
		documentService: documentService,
		projectService:  projectService,
	}
}

func (h *ProposalHandler) CreateProposal(c *gin.Context) {
	// Parse multipart form with 10 MB max memory
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	// Extract form data
	var payload models.Proposal

	// Bind form fields to proposal struct
	if err := c.ShouldBind(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Handle file upload for draft text attachment
	file, header, err := c.Request.FormFile("draft_text_attachment")
	if err == nil {
		// File was uploaded
		defer file.Close()

		// Generate a unique filename
		filename := uuid.New().String() + filepath.Ext(header.Filename)

		// Define upload path (adjust as needed for your application)
		uploadPath := "./uploads/" + filename

		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll("./uploads", 0755); err != nil {
			utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
			return
		}

		// Create the destination file
		out, err := os.Create(uploadPath)
		if err != nil {
			utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create file: "+err.Error())
			return
		}
		defer out.Close()

		// Copy the uploaded file to the destination file
		_, err = io.Copy(out, file)
		if err != nil {
			utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to save file: "+err.Error())
			return
		}

		// Set the file URL in the proposal
		payload.DraftTextAttachmentURL = "/uploads/" + filename
		payload.IsDraftTextAttached = true
	}

	// Set creator ID from authenticated user
	userID, exists := c.Get("userId")
	if !exists {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	payload.CreatedByID = userID.(string)

	// Validate project exists
	projectExists, err := h.projectService.Exists(payload.ProjectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to validate project")
		return
	}

	if !projectExists {
		utilities.ShowMessage(c, http.StatusBadRequest, "Project does not exist")
		return
	}

	// Check if proposal for this project already exists
	exists, err = h.proposalService.Exists(payload.ProjectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to check existing proposals")
		return
	}
	if exists {
		utilities.ShowMessage(c, http.StatusBadRequest, "A proposal for this project already exists")
		return
	}

	err = h.proposalService.Create(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create proposal: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "Proposal created successfully", payload)
}

func (h *ProposalHandler) GetProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	proposal, err := h.proposalService.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Proposal retrieved successfully", proposal)
}

func (h *ProposalHandler) GetProposalByProject(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	proposal, err := h.proposalService.GetByProjectID(projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.Show(c, http.StatusOK, "No proposal found for this project", proposal)
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Proposal retrieved successfully", proposal)
}

func (h *ProposalHandler) ListProposals(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	proposals, total, err := h.proposalService.List(limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to list proposals: "+err.Error())
		return
	}

	pagination := utilities.GeneratePaginationData(limit, page, int(total))
	utilities.Show(c, http.StatusOK, "Proposals retrieved successfully", map[string]any{
		"data":       proposals,
		"pagination": pagination,
	})
}

func (h *ProposalHandler) UpdateProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	// Check if proposal exists
	existing, err := h.proposalService.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	var payload models.Proposal
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

	// Preserve important fields that shouldn't be changed
	payload.ID = id
	payload.CreatedByID = existing.CreatedByID
	payload.CreatedAt = existing.CreatedAt

	err = h.proposalService.Update(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to update proposal: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Proposal updated successfully", payload)
}

func (h *ProposalHandler) UpdatePartialProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	// Check if proposal exists
	_, err = h.proposalService.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid update data: "+err.Error())
		return
	}

	// Prevent updating of critical fields
	delete(updates, "id")
	delete(updates, "created_by_id")
	delete(updates, "created_at")
	delete(updates, "project_id") // Prevent changing the project

	err = h.proposalService.UpdatePartial(id, updates)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to update proposal: "+err.Error())
		return
	}

	// Retrieve the updated proposal
	updated, err := h.proposalService.GetByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve updated proposal: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Proposal updated successfully", updated)
}

func (h *ProposalHandler) DeleteProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	// Check if proposal exists
	_, err = h.proposalService.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	err = h.proposalService.Delete(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to delete proposal: "+err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Proposal deleted successfully")
}

func (h *ProposalHandler) SearchProposals(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Search query is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	proposals, total, err := h.proposalService.Search(query, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to search proposals: "+err.Error())
		return
	}

	pagination := utilities.GeneratePaginationData(limit, page, int(total))
	utilities.Show(c, http.StatusOK, "Search results retrieved", map[string]any{
		"data":       proposals,
		"pagination": pagination,
	})
}

func (h *ProposalHandler) GetProposalsByCreator(c *gin.Context) {
	memberIDStr := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid member ID format")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	proposals, total, err := h.proposalService.GetByCreator(memberID, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposals: "+err.Error())
		return
	}

	pagination := utilities.GeneratePaginationData(limit, page, int(total))
	utilities.Show(c, http.StatusOK, "Proposals retrieved successfully", map[string]any{
		"data":       proposals,
		"pagination": pagination,
	})
}

func (h *ProposalHandler) AddReferencedStandard(c *gin.Context) {
	proposalIDStr := c.Param("id")
	proposalID, err := uuid.Parse(proposalIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	var payload struct {
		DocumentID string `json:"document_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	documentID, err := uuid.Parse(payload.DocumentID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	// Verify both proposal and document exist
	_, err = h.proposalService.GetByID(proposalID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	// Check if document exists - assuming you have a DocumentService
	documentExists, err := h.documentService.Exists(documentID, "", "")
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to validate document")
		return
	}
	if !documentExists {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	err = h.proposalService.AddReferencedStandard(proposalID, documentID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to add referenced standard: "+err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Referenced standard added successfully")
}

func (h *ProposalHandler) RemoveReferencedStandard(c *gin.Context) {
	proposalIDStr := c.Param("id")
	documentIDStr := c.Param("documentId")

	proposalID, err := uuid.Parse(proposalIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	documentID, err := uuid.Parse(documentIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	// Verify proposal exists
	_, err = h.proposalService.GetByID(proposalID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	err = h.proposalService.RemoveReferencedStandard(proposalID, documentID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to remove referenced standard: "+err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Referenced standard removed successfully")
}

func (h *ProposalHandler) TransferProposal(c *gin.Context) {
	proposalIDStr := c.Param("id")
	proposalID, err := uuid.Parse(proposalIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid proposal ID format")
		return
	}

	var payload struct {
		ProjectID string `json:"project_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	newProjectID := payload.ProjectID
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	// Verify proposal exists
	proposal, err := h.proposalService.GetByID(proposalID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utilities.ShowMessage(c, http.StatusNotFound, "Proposal not found")
			return
		}
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve proposal: "+err.Error())
		return
	}

	// Ensure new project is different from current one
	if proposal.ProjectID == newProjectID {
		utilities.ShowMessage(c, http.StatusBadRequest, "New project is the same as current project")
		return
	}

	// Verify new project exists
	projectExists, err := h.projectService.Exists(newProjectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to validate project")
		return
	}
	if !projectExists {
		utilities.ShowMessage(c, http.StatusNotFound, "New project not found")
		return
	}

	// Check if proposal for the new project already exists
	exists, err := h.proposalService.Exists(newProjectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to check existing proposals")
		return
	}

	if exists {
		utilities.ShowMessage(c, http.StatusBadRequest, "A proposal for the target project already exists")
		return
	}

	err = h.proposalService.Transfer(proposalID, newProjectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to transfer proposal: "+err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Proposal transferred successfully")
}
