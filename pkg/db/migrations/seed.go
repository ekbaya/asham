package migrations

import (
	"fmt"
	"log"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeedData struct {
	Stages []models.Stage
}

func GetDefaultSeedData() SeedData {
	return SeedData{
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
	}
}

func SeedDatabase(db *gorm.DB) error {
	seedData := GetDefaultSeedData()

	if err := seedStages(db, seedData.Stages); err != nil {
		return fmt.Errorf("failed to seed stages: %w", err)
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
