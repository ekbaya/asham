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
	return db.AutoMigrate(
		&models.Document{},
		&models.MemberState{},
		&models.NationalStandardBody{},
		&models.Member{},
		&models.ARSOCouncil{},
		&models.JointAdvisoryGroup{},
		&models.StandardsManagementCommittee{},
		&models.TechnicalCommittee{},
		&models.WorkingGroup{},
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
	)
}
