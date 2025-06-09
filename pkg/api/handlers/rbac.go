package handlers

import (
	"net/http"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RbacHandler struct {
	rbacService *services.RbacService
}

func NewRbacHandler(rbacService *services.RbacService) *RbacHandler {
	return &RbacHandler{rbacService: rbacService}
}

func (h *RbacHandler) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure role.Title is uppercase and words are separated by underscores
	processedTitle := utilities.ToUpperUnderscore(role.Title)
	role.Title = processedTitle

	err := h.rbacService.CreateRole(&role)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Role created successfully")
}

func (h *RbacHandler) UpdateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure role.Title is uppercase and words are separated by underscores
	processedTitle := utilities.ToUpperUnderscore(role.Title)
	role.Title = processedTitle

	err := h.rbacService.UpdateRole(&role)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Role created successfully")
}

func (h *RbacHandler) ListRoles(c *gin.Context) {
	roles, err := h.rbacService.ListRoles()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "roles", roles)
}

func (h *RbacHandler) GetRoleByID(c *gin.Context) {
	role, err := h.rbacService.GetRoleByID(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "role", role)
}

func (h *RbacHandler) DeleteRole(c *gin.Context) {
	err := h.rbacService.DeleteRole(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "role has been deleted")
}

func (h *RbacHandler) CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure permission.Title is uppercase and words are separated by underscores
	processedTitle := utilities.ToUpperUnderscore(permission.Title)
	permission.Title = processedTitle

	err := h.rbacService.CreatePermission(&permission)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Permission created successfully")
}

func (h *RbacHandler) ListPermissions(c *gin.Context) {
	perms, err := h.rbacService.ListPermissions()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "permissions", perms)
}

func (h *RbacHandler) AssignRoleToMember(c *gin.Context) {
	memberID := c.Param("member_id")
	roleID := c.Param("role_id")

	if err := h.rbacService.AssignRoleToMember(memberID, roleID); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Role assigned to member successfully")
}

func (h *RbacHandler) ListMemberRoles(c *gin.Context) {
	memberID := c.Param("member_id")

	roles, err := h.rbacService.ListRolesForMember(memberID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "roles", roles)
}

func (h *RbacHandler) DeletePermission(c *gin.Context) {
	err := h.rbacService.DeletePermission(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "role has been deleted")
}
