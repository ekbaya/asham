package migrations

import (
	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
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
	)
}
