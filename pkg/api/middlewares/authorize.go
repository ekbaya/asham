package middleware

import (
	"fmt"
	"net/http"

	"github.com/ekbaya/asham/pkg/config"
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
		globalConfig := config.GetConfig()
		fmt.Println("ENV CHECK >>>", globalConfig.Environment)

		if globalConfig.Environment == "dev" {
			c.Next()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			fmt.Println("DEBUG >>> User not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		member, err := service.Account(userID.(string))
		if err != nil {
			fmt.Println("DEBUG >>> Failed to get user account:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		permSlug, err := service.GetPermissionSlug(c.Request.Method, c.FullPath())
		fmt.Printf("PERM SLUG >>> method=%s, path=%s, slug=%s, err=%v\n", c.Request.Method, c.FullPath(), permSlug, err)
		if err != nil || permSlug == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission not mapped"})
			return
		}

		if hasPermission(member.Roles, permSlug) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("This action is forbidden. Permission: %s is required", permSlug)})
	}
}

func hasPermission(roles []models.Role, required string) bool {
	fmt.Println("CHECKING PERMISSION >>> Required:", required)
	permSet := make(map[string]bool)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			permSet[perm.Slug] = true
		}
	}
	return permSet[required]
}
