package app

import (
	"fmt"
	"strings"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/gin-gonic/gin"
)

func registerPermissionsAndResources(router *gin.Engine, permissionService *services.RbacService, resourceService *services.PermissionResourceService) error {
	for _, r := range router.Routes() {
		method := r.Method
		path := r.Path

		// Generate slug based on convention: METHOD + cleaned path
		slug := generatePermissionSlug(method, path)

		// Create or get permission
		permission := models.Permission{
			Title:       strings.ToUpper(slug),
			Description: fmt.Sprintf("Allows %s %s", method, path),
			Action:      strings.ToLower(method),
			Resource:    extractResource(path),
		}

		// Create if not exists
		perm, err := permissionService.CreateIfNotExists(permission)
		if err != nil {
			fmt.Printf("failed to create permission: %v\n", err)
			continue
		}

		// Create method+path mapping to permission
		err = resourceService.CreateMapping(method, path, perm.ID)
		if err != nil {
			fmt.Printf("failed to create route-permission mapping: %v\n", err)
		}
	}
	return nil
}

func generatePermissionSlug(method, path string) string {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	var cleanedParts []string
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			continue
		}
		cleanedParts = append(cleanedParts, part)
	}
	resource := strings.Join(cleanedParts, "_")
	return fmt.Sprintf("%s_%s", strings.ToLower(method), resource)
}

func extractResource(path string) string {
	// You can enhance this logic later
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if !strings.HasPrefix(p, ":") && p != "" {
			return p
		}
	}
	return "unknown"
}
