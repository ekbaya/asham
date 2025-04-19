package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	documentService services.DocumentService
}

func NewDocumentHandler(documentService services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(100 << 20); err != nil { // 100 MB max
		utilities.ShowMessage(c, http.StatusBadRequest, "Unable to parse form: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	project := c.PostForm("project")
	docType := c.PostForm("type")

	if project == "" || docType == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Project and type are required fields")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	// Add validation for ARS document type - must be PDF
	if docType == "ARS" {
		extension := strings.ToLower(filepath.Ext(header.Filename))
		if extension != ".pdf" {
			utilities.ShowMessage(c, http.StatusBadRequest, "ARS standard must be a PDF file")
			return
		}
	}

	assetsDir := "../assets/documents"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create assets directory: "+err.Error())
		return
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filepath := filepath.Join(assetsDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create destination file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to save file: "+err.Error())
		return
	}

	// Set the file URL in the document
	// For local storage, we'll use a relative path
	// This we will change to an S3 URL later
	fileURl := "/" + filepath // prepend with slash for URL format

	err = h.documentService.UpdateProjectDoc(project, docType, fileURl, userID.(string))
	if err != nil {
		os.Remove(filepath)
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, fmt.Sprintf("%s added successfully", docType))
}

// CreateDocument handles the creation of a new document with file upload
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		utilities.ShowMessage(c, http.StatusBadRequest, "Unable to parse form: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var payload models.Document
	payload.Title = c.PostForm("title")
	payload.Reference = c.PostForm("reference")
	payload.Description = c.PostForm("description")

	payload.CreatedByID = userID.(string)

	if payload.Title == "" || payload.Reference == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Title and reference are required fields")
		return
	}

	exists, err := h.documentService.Exists(uuid.Nil, payload.Reference, payload.Title)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	if exists {
		utilities.ShowMessage(c, http.StatusConflict, "Document with the same reference or title already exists")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	assetsDir := "../assets/documents"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create assets directory: "+err.Error())
		return
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filepath := filepath.Join(assetsDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create destination file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to save file: "+err.Error())
		return
	}

	// Set the file URL in the document
	// For local storage, we'll use a relative path
	// This we will change to an S3 URL later
	payload.FileURL = "/" + filepath // prepend with slash for URL format

	err = h.documentService.Create(&payload)
	if err != nil {
		os.Remove(filepath)
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "Document added successfully", payload)
}

// GetDocumentByID retrieves a document by its ID
func (h *DocumentHandler) GetDocumentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	document, err := h.documentService.GetByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if document == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	utilities.Show(c, http.StatusOK, "document", document)
}

// GetDocumentByReference retrieves a document by its reference
func (h *DocumentHandler) GetDocumentByReference(c *gin.Context) {
	reference := c.Param("reference")

	document, err := h.documentService.GetByReference(reference)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if document == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	utilities.Show(c, http.StatusOK, "document", document)
}

// GetDocumentByTitle retrieves a document by its title
func (h *DocumentHandler) GetDocumentByTitle(c *gin.Context) {
	title := c.Param("title")

	document, err := h.documentService.GetByTitle(title)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if document == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	utilities.Show(c, http.StatusOK, "document", document)
}

// UpdateDocument updates an existing document
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	// Check if document exists
	existingDoc, err := h.documentService.GetByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingDoc == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	var payload models.Document
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

	// Set the ID to ensure we update the correct document
	payload.ID = id

	err = h.documentService.Update(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Document updated successfully")
}

// UpdateDocumentPartial updates specific fields of a document
func (h *DocumentHandler) UpdateDocumentPartial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	// Check if document exists
	existingDoc, err := h.documentService.GetByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingDoc == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.documentService.UpdatePartial(id, updates)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Document partially updated successfully")
}

// UpdateFileURL updates only the file URL of a document
func (h *DocumentHandler) UpdateFileURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	// Check if document exists
	existingDoc, err := h.documentService.GetByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingDoc == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	var payload struct {
		FileURL string `json:"file_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.documentService.UpdateFileURL(id, payload.FileURL)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Document file URL updated successfully")
}

// DeleteDocument deletes a document by its ID
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid document ID format")
		return
	}

	// Check if document exists
	existingDoc, err := h.documentService.GetByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingDoc == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Document not found")
		return
	}

	err = h.documentService.Delete(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Document deleted successfully")
}

// ListDocuments retrieves a paginated list of documents
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	limit := utilities.IntQueryParam(c, "limit", 10)
	page := utilities.IntQueryParam(c, "page", 1)
	offset := (page - 1) * limit

	documents, total, err := h.documentService.List(limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := gin.H{
		"documents": documents,
		"pagination": gin.H{
			"total": total,
			"limit": limit,
			"page":  page,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utilities.Show(c, http.StatusOK, "data", response)
}

// SearchDocuments searches for documents based on a query
func (h *DocumentHandler) SearchDocuments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Search query is required")
		return
	}

	limit := utilities.IntQueryParam(c, "limit", 10)
	page := utilities.IntQueryParam(c, "page", 1)
	offset := (page - 1) * limit

	documents, total, err := h.documentService.Search(query, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := gin.H{
		"documents": documents,
		"pagination": gin.H{
			"total": total,
			"limit": limit,
			"page":  page,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utilities.Show(c, http.StatusOK, "data", response)
}

// GetDocumentsByDateRange retrieves documents created within a specified date range
func (h *DocumentHandler) GetDocumentsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Both start_date and end_date are required")
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	// Add a day to end date to include the entire end date
	endDate = endDate.Add(24 * time.Hour)

	documents, err := h.documentService.GetDocumentsCreatedBetween(startDate, endDate)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "documents", documents)
}

// CountDocuments returns the total count of documents
func (h *DocumentHandler) CountDocuments(c *gin.Context) {
	count, err := h.documentService.CountAll()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "count", count)
}

func (h *DocumentHandler) ProjectDocuments(c *gin.Context) {
	docs, err := h.documentService.ProjectDocuments(c.Param("projectId"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "docs", docs)
}

func (h *DocumentHandler) UploadRelatedDocument(c *gin.Context) {
	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		utilities.ShowMessage(c, http.StatusBadRequest, "Unable to parse form: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	project := c.PostForm("project")
	docTitle := c.PostForm("title")
	docDesc := c.PostForm("description")
	reference := c.PostForm("reference")

	if project == "" || docTitle == "" || reference == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Project, title and reference are required fields")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	assetsDir := "../assets/documents"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create assets directory: "+err.Error())
		return
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filepath := filepath.Join(assetsDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create destination file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to save file: "+err.Error())
		return
	}

	// Set the file URL in the document
	// For local storage, we'll use a relative path
	// This we will change to an S3 URL later
	fileURl := "/" + filepath // prepend with slash for URL format

	err = h.documentService.UpdateProjectRelatedDoc(project, docTitle, reference, docDesc, fileURl, userID.(string))
	if err != nil {
		os.Remove(filepath)
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, fmt.Sprintf("%s added successfully", docTitle))
}

func (h *DocumentHandler) UploadStandard(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(100 << 20); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Unable to parse form: "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var payload models.Document
	payload.Title = c.PostForm("title")
	payload.Reference = c.PostForm("reference")
	payload.Description = c.PostForm("description")
	payload.CreatedByID = userID.(string)

	sector := c.PostForm("sector")
	language := c.PostForm("language")
	description := c.PostForm("description")
	tc := c.PostForm("tc")

	project := models.Project{
		MemberID:             payload.CreatedByID,
		ProjectSectorID:      &sector,
		Reference:            payload.Reference,
		Title:                payload.Title,
		Language:             language,
		Description:          description,
		TechnicalCommitteeID: tc,
	}

	if payload.Title == "" || payload.Reference == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Title and reference are required fields")
		return
	}

	exists, err := h.documentService.Exists(uuid.Nil, payload.Reference, payload.Title)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	if exists {
		utilities.ShowMessage(c, http.StatusConflict, "Document with the same reference or title already exists")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	assetsDir := "../assets/documents"
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create assets directory: "+err.Error())
		return
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filepath := filepath.Join(assetsDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create destination file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to save file: "+err.Error())
		return
	}

	// Set the file URL in the document
	// For local storage, we'll use a relative path
	// This we will change to an S3 URL later
	payload.FileURL = "/" + filepath // prepend with slash for URL format

	err = h.documentService.UploadStandard(&payload, &project)
	if err != nil {
		os.Remove(filepath)
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "Document added successfully", payload)
}
