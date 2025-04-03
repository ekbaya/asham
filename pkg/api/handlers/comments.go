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

// CommentHandler struct
type CommentHandler struct {
	commentService *services.CommentService
}

// NewCommentHandler initializes a new CommentHandler
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var payload models.CommentObservation
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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	payload.NationalSecretaryID = userID.(string)

	err := h.commentService.Create(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Comment added successfully")
}

func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := h.commentService.GetByID(commentID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "comment", comment)
}

func (h *CommentHandler) GetAllComments(c *gin.Context) {
	comments, err := h.commentService.GetAll()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "comments", comments)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("comment_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	var payload models.CommentObservation
	if err := c.ShouldBindJSON(&payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	payload.ID = commentID // Ensure ID matches request parameter
	err = h.commentService.Update(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Comment updated successfully")
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	err = h.commentService.Delete(commentID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Comment deleted successfully")
}

func (h *CommentHandler) GetCommentsByProjectID(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	comments, err := h.commentService.GetByProjectID(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "comments", comments)
}
