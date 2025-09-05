package migrations

import (
	"fmt"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Ensure UUID extension is created
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}
	
	// First, create tables without foreign key constraints
	// We'll use DisableForeignKeyConstraintWhenMigrating
	originalConfig := db.Config
	db.Config.DisableForeignKeyConstraintWhenMigrating = true
	
	err := db.AutoMigrate(
		&models.Sector{},
		&models.Document{},
		&models.MemberState{},
		&models.Permission{},
		&models.Role{},
		&models.Member{},
		&models.NationalStandardBody{},
		&models.ARSOCouncil{},
		&models.JointAdvisoryGroup{},
		&models.StandardsManagementCommittee{},
		&models.TechnicalCommittee{},
		&models.WorkingGroup{},
		&models.EditingCommittee{},
		&models.TaskForce{},
		&models.SubCommittee{},
		&models.SpecializedCommittee{},
		&models.JointTechnicalCommittee{},
		&models.ProjectDuration{},
		&models.Timeframe{},
		&models.Stage{},
		&models.Project{},
		&models.ProjectStageHistory{},
		&models.Proposal{},
		&models.Acceptance{},
		&models.NSBResponse{},
		&models.CommentObservation{},
		&models.NSBResponseStatusChange{},
		&models.DARS{},
		&models.NationalConsultation{},
		&models.Balloting{},
		&models.Vote{},
		&models.Meeting{},
		&models.User{},
		&models.Standard{},
		&models.StandardVersion{},
		&models.StandardAuditLog{},
		&models.ResourcePermission{},
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.NotificationPreference{},
		&models.NotificationHistory{},
	)
	
	// Restore original config
	db.Config = originalConfig
	
	return err
}
