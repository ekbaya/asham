package migrations

import (
	"fmt"
	"log"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeedData struct {
	AdminUser       models.Member
	Stages          []models.Stage
	Sectors         []string
	Roles           []models.Role
	Permissions     []models.Permission
	RolePermissions map[string][]string
}

func GetDefaultSeedData() SeedData {
	return SeedData{
		AdminUser: models.Member{
			FirstName: "ARSO",
			LastName:  "ARSO",
			Phone:     "254712345678",
			Email:     "support@arso.com",
		},
		Sectors: []string{"Health", "IT and Related Services", "Management and Related Services", "Safety, Security and Risk Services", "Transport", "Energy", "Diversity and Inclusion", "Environment and Sustainability", "Food and Agriculture", "Building and Construction", "Engineering"},
		Stages: []models.Stage{
			{
				ID:           uuid.New(),
				Number:       0,
				Name:         "Preliminary stage",
				DocumentName: "Preliminary Work Item",
				Abbreviation: "PWI",
				CreatedAt:    time.Now(),
			},
			{
				ID:           uuid.New(),
				Number:       1,
				Name:         "Proposal stage",
				DocumentName: "New Work Item Proposal",
				Abbreviation: "NWIP",
				CreatedAt:    time.Now(),
				Timeframe: &models.Timeframe{
					ID:          uuid.New(),
					Standard:    &models.ProjectDuration{ID: uuid.New(), Min: 3 * 30, Max: 5 * 30}, // 3 to 5 months
					IS:          &models.ProjectDuration{ID: uuid.New(), Min: 30},                  // 1 month
					Emergency:   &models.ProjectDuration{ID: uuid.New(), Min: 21},                  // 21 days
					Description: "Standard: 3 - 5 months; IS: 1 month; Emergency: 21 days",
				},
			},
			{
				ID:           uuid.New(),
				Number:       2,
				Name:         "Preparatory stage",
				DocumentName: "Working Draft(s)",
				Abbreviation: "WD",
				CreatedAt:    time.Now(),
				Timeframe: &models.Timeframe{
					ID:          uuid.New(),
					Standard:    &models.ProjectDuration{ID: uuid.New(), Min: 2 * 30}, // 2 months
					Description: "Standard: 2 months",
				},
			},
			{
				ID:           uuid.New(),
				Number:       3,
				Name:         "Committee stage",
				DocumentName: "Committee Draft(s)",
				Abbreviation: "CD",
				CreatedAt:    time.Now(),
				Timeframe: &models.Timeframe{
					ID:          uuid.New(),
					Standard:    &models.ProjectDuration{ID: uuid.New(), Min: 6 * 30},
					Emergency:   &models.ProjectDuration{ID: uuid.New(), Min: 15}, // 15 days
					Description: "Standard: 6 months, Emergency: 15 days",
				},
			},
			{
				ID:           uuid.New(),
				Number:       4,
				Name:         "Enquiry stage",
				DocumentName: "Draft African Standard",
				Abbreviation: "DARS",
				CreatedAt:    time.Now(),
				Timeframe: &models.Timeframe{
					ID:          uuid.New(),
					Standard:    &models.ProjectDuration{ID: uuid.New(), Min: 4 * 30}, // 4 months
					IS:          &models.ProjectDuration{ID: uuid.New(), Min: 2 * 30}, // 2 months
					Emergency:   &models.ProjectDuration{ID: uuid.New(), Min: 30},     // 30 days
					Description: "Standard: 4 months; IS: 2 month; Emergency: 30 days",
				},
			},
			{
				ID:           uuid.New(),
				Number:       5,
				Name:         "Ballot stage",
				DocumentName: "Final Draft African Standard",
				Abbreviation: "FDARS",
				CreatedAt:    time.Now(),
				Timeframe: &models.Timeframe{
					ID:          uuid.New(),
					Standard:    &models.ProjectDuration{ID: uuid.New(), Min: 30}, // 1 months
					IS:          &models.ProjectDuration{ID: uuid.New(), Min: 30}, // 1 months
					Emergency:   &models.ProjectDuration{ID: uuid.New(), Min: 6},  // 6 days
					Description: "Standard: 1 month; IS: 1 month; Emergency: 6 days",
				},
			},
			{
				ID:           uuid.New(),
				Number:       6,
				Name:         "Approval stage",
				DocumentName: "Final Draft African Standard",
				Abbreviation: "FDARS",
				CreatedAt:    time.Now(),
				Timeframe: &models.Timeframe{
					ID:          uuid.New(),
					Standard:    &models.ProjectDuration{ID: uuid.New(), Min: 3 * 30}, // 3 months
					IS:          &models.ProjectDuration{ID: uuid.New(), Min: 3 * 30}, // 3 months
					Emergency:   &models.ProjectDuration{ID: uuid.New(), Min: 15},     // 15 days
					Description: "Standard: 3 month; IS: 3 month; Emergency: 15 days",
				},
			},
		},
		Roles: []models.Role{
			{
				Title:       "TC_SECRETARIAT",
				Description: "Technical Committee Secretariat",
			},
			{
				Title:       "SC_MEMBER",
				Description: "Sub Committee Member",
			},
			{
				Title:       "SMC_MEMBER",
				Description: "Standards Management Committee Member",
			},
			{
				Title:       "NSB_MEMBER",
				Description: "National Standards Body Member",
			},
			{
				Title:       "PROJECT_LEADER",
				Description: "Leader of specific standardization projects",
			},
			{
				Title:       "WORKING_GROUP_MEMBER",
				Description: "Member of Working Groups",
			},
			{
				Title:       "NATIONAL_EXPERT",
				Description: "Expert at national level",
			},
			{
				Title:       "ARSO_CENTRAL_SECRETARIAT",
				Description: "ARSO Central Secretariat Member",
			},
		},
		Permissions: []models.Permission{
			// Preliminary Stage Permissions
			{
				Title:       "EVALUATE_PWI",
				Description: "Can evaluate Preliminary Work Items",
				Resource:    "preliminary_work_item",
				Action:      "evaluate",
			},
			{
				Title:       "REVIEW_PWI",
				Description: "Can review Preliminary Work Items",
				Resource:    "preliminary_work_item",
				Action:      "review",
			},

			// Proposal Stage Permissions
			{
				Title:       "CREATE_NWIP",
				Description: "Can create New Work Item Proposals",
				Resource:    "new_work_item_proposal",
				Action:      "create",
			},
			{
				Title:       "VOTE_NWIP",
				Description: "Can vote on New Work Item Proposals",
				Resource:    "new_work_item_proposal",
				Action:      "vote",
			},
			{
				Title:       "COMPILE_POSITIONS",
				Description: "Can compile national positions",
				Resource:    "national_position",
				Action:      "compile",
			},

			// Preparatory Stage Permissions
			{
				Title:       "DEVELOP_WORKING_DRAFT",
				Description: "Can develop Working Drafts",
				Resource:    "working_draft",
				Action:      "develop",
			},
			{
				Title:       "INVITE_EXPERTS",
				Description: "Can invite expert assistance",
				Resource:    "expert_assistance",
				Action:      "invite",
			},
			{
				Title:       "ASSIGN_PROJECT_LEADER",
				Description: "Can assign Project Leaders",
				Resource:    "project_leader",
				Action:      "assign",
			},

			// Committee Stage Permissions
			{
				Title:       "REVIEW_COMMITTEE_DRAFT",
				Description: "Can review Committee Drafts",
				Resource:    "committee_draft",
				Action:      "review",
			},
			{
				Title:       "SUBMIT_COMMENTS",
				Description: "Can submit comments on drafts",
				Resource:    "draft_comments",
				Action:      "submit",
			},
			{
				Title:       "REGISTER_ENQUIRY_STAGE",
				Description: "Can register drafts for Enquiry Stage",
				Resource:    "enquiry_stage",
				Action:      "register",
			},

			// Enquiry Stage Permissions
			{
				Title:       "GENERATE_DARS",
				Description: "Can generate Draft African Standards",
				Resource:    "draft_african_standard",
				Action:      "generate",
			},
			{
				Title:       "PUBLIC_REVIEW",
				Description: "Can participate in public review",
				Resource:    "public_review",
				Action:      "review",
			},
		},
		RolePermissions: map[string][]string{
			"TC_SECRETARIAT": {
				"EVALUATE_PWI",
				"REVIEW_PWI",
				"CREATE_NWIP",
				"ASSIGN_PROJECT_LEADER",
				"REGISTER_ENQUIRY_STAGE",
			},
			"SMC_MEMBER": {
				"EVALUATE_PWI",
				"VOTE_NWIP",
				"COMPILE_POSITIONS",
			},
			"NSB_MEMBER": {
				"SUBMIT_COMMENTS",
				"VOTE_NWIP",
				"PUBLIC_REVIEW",
			},
			"PROJECT_LEADER": {
				"DEVELOP_WORKING_DRAFT",
				"INVITE_EXPERTS",
			},
			"WORKING_GROUP_MEMBER": {
				"DEVELOP_WORKING_DRAFT",
				"REVIEW_COMMITTEE_DRAFT",
				"SUBMIT_COMMENTS",
			},
			"NATIONAL_EXPERT": {
				"REVIEW_COMMITTEE_DRAFT",
				"SUBMIT_COMMENTS",
			},
			"ARSO_CENTRAL_SECRETARIAT": {
				"GENERATE_DARS",
				"REGISTER_ENQUIRY_STAGE",
			},
		},
	}
}

func SeedDatabase(db *gorm.DB) error {
	seedData := GetDefaultSeedData()

	if err := seedAdminUser(db, &seedData.AdminUser); err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	if err := seedStages(db, seedData.Stages); err != nil {
		return fmt.Errorf("failed to seed stages: %w", err)
	}

	if err := seedSectors(db, seedData.Sectors); err != nil {
		return fmt.Errorf("failed to seed sectors: %w", err)
	}

	if err := seedRoles(db, seedData.Roles); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	if err := seedPermissions(db, seedData.Permissions); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	if err := seedRolePermissions(db, seedData.RolePermissions); err != nil {
		return fmt.Errorf("failed to seed role permissions: %w", err)
	}

	return nil
}

func seedSectors(db *gorm.DB, sectorsList []string) error {
	var count int64
	if err := db.Model(&models.Sector{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check sectors: %w", err)
	}

	if count == 0 {
		statuses := make([]models.Sector, len(sectorsList))
		for i, status := range sectorsList {
			statuses[i] = models.Sector{Title: status}
		}

		if err := db.Create(&statuses).Error; err != nil {
			return fmt.Errorf("failed to create payment statuses: %w", err)
		}
		log.Printf("Created %d payment statuses", len(sectorsList))
	}

	return nil
}

func seedAdminUser(db *gorm.DB, adminData *models.Member) error {
	var count int64
	if err := db.Model(models.Member{}).Where("email = ?", adminData.Email).Count(&count).Error; err != nil {
		log.Printf("Error counting user: %v", err)
		return err
	}

	if condition := count == 0; condition {
		// Generate a single user ID that will be used consistently
		userID := uuid.New()
		log.Printf("Generated user ID: %s", userID)

		// Hash password
		hashedPassword, err := utilities.HashPassword("secret")
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		// Prepare user data
		adminUser := models.Member{
			ID:             userID,
			FirstName:      adminData.FirstName,
			LastName:       adminData.LastName,
			Phone:          adminData.Phone,
			Email:          adminData.Email,
			HashedPassword: hashedPassword,
		}

		if err := db.FirstOrCreate(&adminUser).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			return err
		}
	}
	return nil
}

func seedStages(db *gorm.DB, stageList []models.Stage) error {
	var count int64
	if err := db.Model(&models.Stage{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check stages: %w", err)
	}

	if count == 0 {
		// Create stages with their associated timeframes and ProjectDurations in a transaction
		return db.Transaction(func(tx *gorm.DB) error {
			for _, stage := range stageList {
				// First, create any non-nil ProjectDurations
				if stage.Timeframe != nil {
					if stage.Timeframe.Standard != nil {
						if err := tx.Create(stage.Timeframe.Standard).Error; err != nil {
							return fmt.Errorf("failed to create standard ProjectDuration: %w", err)
						}
					}
					if stage.Timeframe.IS != nil {
						if err := tx.Create(stage.Timeframe.IS).Error; err != nil {
							return fmt.Errorf("failed to create IS ProjectDuration: %w", err)
						}
					}
					if stage.Timeframe.Emergency != nil {
						if err := tx.Create(stage.Timeframe.Emergency).Error; err != nil {
							return fmt.Errorf("failed to create emergency ProjectDuration: %w", err)
						}
					}

					// Then create the timeframe
					if err := tx.Create(stage.Timeframe).Error; err != nil {
						return fmt.Errorf("failed to create timeframe: %w", err)
					}
				}

				// Finally, create the stage
				if err := tx.Create(&stage).Error; err != nil {
					return fmt.Errorf("failed to create stage: %w", err)
				}
			}

			log.Printf("Created %d stages", len(stageList))
			return nil
		})
	}

	return nil
}

func seedRoles(db *gorm.DB, roles []models.Role) error {
	for _, role := range roles {
		var count int64
		if err := db.Model(&models.Role{}).Where("title = ?", role.Title).Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			role.ID = uuid.New()
			if err := db.Create(&role).Error; err != nil {
				return err
			}
			log.Printf("Created role: %s", role.Title)
		}
	}
	return nil
}

func seedPermissions(db *gorm.DB, permissions []models.Permission) error {
	for _, permission := range permissions {
		var count int64
		if err := db.Model(&models.Permission{}).Where("title = ?", permission.Title).Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			permission.ID = uuid.New()
			if err := db.Create(&permission).Error; err != nil {
				return err
			}
			log.Printf("Created permission: %s", permission.Title)
		}
	}
	return nil
}

func seedRolePermissions(db *gorm.DB, rolePermissions map[string][]string) error {
	// Create a join table if it doesn't exist
	if err := db.Table("role_permissions").AutoMigrate(&struct {
		RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
		PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
	}{}); err != nil {
		return fmt.Errorf("failed to ensure role_permissions table exists: %w", err)
	}

	for roleName, permissionNames := range rolePermissions {
		// Find the role
		var role models.Role
		if err := db.Where("title = ?", roleName).First(&role).Error; err != nil {
			log.Printf("Warning: Role %s not found: %v", roleName, err)
			continue
		}

		// Process each permission for the role
		for _, permissionName := range permissionNames {
			var permission models.Permission
			if err := db.Where("title = ?", permissionName).First(&permission).Error; err != nil {
				log.Printf("Warning: Permission %s not found: %v", permissionName, err)
				continue
			}

			// Check if relationship already exists
			var exists bool
			err := db.Table("role_permissions").
				Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
				Select("count(*) > 0").
				Scan(&exists).Error

			if err != nil {
				return fmt.Errorf("failed to check role-permission relationship: %w", err)
			}

			if !exists {
				// Create the relationship using SQL to ensure proper insertion
				err = db.Exec(`
                    INSERT INTO role_permissions (role_id, permission_id)
                    VALUES (?, ?)
                `, role.ID, permission.ID).Error

				if err != nil {
					return fmt.Errorf("failed to create role-permission relationship: %w", err)
				}

				log.Printf("Added permission %s to role %s", permissionName, roleName)
			} else {
				log.Printf("Permission %s already exists for role %s", permissionName, roleName)
			}
		}
	}

	// Verify relationships were created successfully
	for roleName, permissionNames := range rolePermissions {
		var role models.Role
		if err := db.Preload("Permissions").Where("name = ?", roleName).First(&role).Error; err != nil {
			log.Printf("Warning: Could not verify role %s: %v", roleName, err)
			continue
		}

		if len(role.Permissions) != len(permissionNames) {
			log.Printf("Warning: Role %s has %d permissions, expected %d",
				roleName, len(role.Permissions), len(permissionNames))
		}
	}

	return nil
}
