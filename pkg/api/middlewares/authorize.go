package middleware

import (
	"net/http"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/gin-gonic/gin"
)

func Authorize(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		member := user.(*models.Member)

		for _, required := range requiredPermissions {
			if hasPermission(member.Roles, required) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

func DynamicAuthorize(service *services.PermissionResourceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		member := user.(*models.Member)

		permSlug, err := service.GetPermissionSlug(c.Request.Method, c.FullPath())
		if err != nil || permSlug == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission not mapped"})
			return
		}

		if hasPermission(member.Roles, permSlug) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

func hasPermission(roles []models.Role, required string) bool {
	permSet := make(map[string]bool)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			permSet[perm.Slug] = true
		}
	}
	return permSet[required]
}
