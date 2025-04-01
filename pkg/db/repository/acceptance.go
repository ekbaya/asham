package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AcceptanceRepository struct {
	db *gorm.DB
}

func NewAcceptanceRepository(db *gorm.DB) *AcceptanceRepository {
	return &AcceptanceRepository{db: db}
}

func (r *AcceptanceRepository) CreateNSBResponse(response *models.NSBResponse) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var acceptance models.Acceptance

		// Check if Acceptance exists for the given project
		if err := tx.Where("project_id = ?", response.Project).First(&acceptance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create a new Acceptance if it does not exist
				acceptance = models.Acceptance{
					ID:         uuid.New(),
					ProjectID:  response.ProjectID,
					IsApproved: false,
					CreatedAt:  time.Now(),
				}
				if err := tx.Create(&acceptance).Error; err != nil {
					return err
				}
			} else {
				return err // Return unexpected errors
			}
		}

		var member models.Member
		if err := tx.Where("id = ?", response.ResponderID).First(&member).Error; err != nil {
			return err
		}

		// Attach NSBResponse to the existing acceptance
		response.AcceptanceID = acceptance.ID
		nsb := member.NationalStandardBodyID
		response.RespondingNSBID = *nsb
		if err := tx.Create(response).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *AcceptanceRepository) GetNSBResponse(id string) (*models.NSBResponse, error) {
	var response models.NSBResponse
	if err := r.db.Where("id = ?", id).First(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *AcceptanceRepository) GetNSBResponsesByProjectID(projectID string) ([]models.NSBResponse, error) {
	var NSBResponses []models.NSBResponse
	if err := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&NSBResponses).Error; err != nil {
		return nil, err
	}
	return NSBResponses, nil
}

func (r *AcceptanceRepository) UpdateNSBResponse(response *models.NSBResponse) error {
	return r.db.Save(response).Error
}

func (r *AcceptanceRepository) DeleteNSBResponse(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.NSBResponse{}).Error
}

func (r *AcceptanceRepository) GetAcceptance(id string) (*models.Acceptance, error) {
	var acceptance models.Acceptance
	if err := r.db.Where("id = ?", id).First(&acceptance).Error; err != nil {
		return nil, err
	}
	return &acceptance, nil
}

func (r *AcceptanceRepository) GetAcceptanceByProjectID(projectID string) (*models.Acceptance, error) {
	var acceptance models.Acceptance
	if err := r.db.Where("project_id = ?", projectID).First(&acceptance).Error; err != nil {
		return nil, err
	}
	return &acceptance, nil
}

func (r *AcceptanceRepository) GetAcceptances() (*[]models.Acceptance, error) {
	var acceptances []models.Acceptance
	if err := r.db.Preload(clause.Associations).Order("created_at DESC").Find(&acceptances).Error; err != nil {
		return nil, err
	}
	return &acceptances, nil
}

func (r *AcceptanceRepository) UpdateAcceptance(acceptance *models.Acceptance) error {
	return r.db.Save(acceptance).Error
}

func (r *AcceptanceRepository) GetAcceptanceWithResponses(id string) (*models.Acceptance, error) {
	var acceptance models.Acceptance
	if err := r.db.Where("id = ?", id).Preload("acceptances").First(&acceptance).Error; err != nil {
		return nil, err
	}
	return &acceptance, nil
}

func (r *AcceptanceRepository) CountNSBResponsesByType(projectID string) (map[models.Response]int, error) {
	var results []struct {
		Response models.Response
		Count    int
	}

	if err := r.db.Model(&models.NSBResponse{}).
		Select("response, count(*) as count").
		Where("project_id = ?", projectID).
		Group("response").
		Find(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[models.Response]int)
	for _, result := range results {
		counts[result.Response] = result.Count
	}

	return counts, nil
}

func (r *AcceptanceRepository) CalculateNSBResponseStats(projectID string) error {
	// Count responses by type
	counts, err := r.CountNSBResponsesByType(projectID)
	if err != nil {
		return err
	}

	// Initialize counters
	totalResponses := 0
	agreementCount := 0
	disagreementCount := 0
	abstentionCount := 0

	// Process response counts
	for response, count := range counts {
		totalResponses += count

		switch response {
		case models.ResponseAgreeAdvance,
			models.ResponseAgreeAcceptWorkingDraft,
			models.ResponseAgreeCirculateCD,
			models.ResponseAgreeCirculateDARF:
			agreementCount += count
		case models.ResponseNoAgreement:
			disagreementCount += count
		case models.ResponseAbstention:
			abstentionCount += count
		}
	}

	// Perform raw SQL update
	query := `
		UPDATE acceptances
		SET total_responses = ?, agreement_count = ?, disagreement_count = ?, abstention_count = ?
		WHERE project_id = ?
	`
	result := r.db.Exec(query, totalResponses, agreementCount, disagreementCount, abstentionCount, projectID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *AcceptanceRepository) SetAcceptanceApproval(id string, approved bool) error {
	var acceptance models.Acceptance
	if err := r.db.Where("id = ?", id).First(&acceptance).Error; err != nil {
		return err
	}

	acceptance.IsApproved = approved
	if approved {
		now := time.Now()
		acceptance.SMCApprovalDate = &now
	} else {
		acceptance.SMCApprovalDate = &time.Time{} // Zero value
	}

	return r.db.Save(&acceptance).Error
}

func (r *AcceptanceRepository) GetAcceptanceResults(id string) (*models.AcceptanceResults, error) {
	// Parse the UUID from string
	acceptanceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	// Find the Acceptance by ID
	var acceptance models.Acceptance
	if err := r.db.Preload("Submissions").Preload("Submissions.RespondingNSB").First(&acceptance, acceptanceID).Error; err != nil {
		return nil, err
	}

	// Initialize results structure
	results := &models.AcceptanceResults{
		AcceptanceID:           acceptance.ID,
		ProjectID:              acceptance.ProjectID,
		IndividualNSBResponses: make([]models.IndividualNSBResponse, 0),
		Totals:                 models.ResponseTotals{},
	}

	// Process each NSB response
	for _, response := range *acceptance.Submissions {
		nsbResponse := models.IndividualNSBResponse{
			NSB:                    response.RespondingNSB.Name,
			FeasibleYes:            false,
			FeasibleNo:             false,
			Abstention:             false,
			AcceptedAsNWIP:         "N",
			AcceptedForProgressing: "N",
			AcceptedAsWD:           "N",
			AcceptedAsCD:           "N",
			AcceptedAsDARF:         "N",
			CommentsEnclosed:       response.Comments != "",
			Participation:          response.IsCommittedToParticipate,
		}

		// Determine feasibility response
		switch response.Response {
		case models.ResponseAgreeAdvance, models.ResponseAgreeAcceptWorkingDraft,
			models.ResponseAgreeCirculateCD, models.ResponseAgreeCirculateDARF:
			nsbResponse.FeasibleYes = true
			results.Totals.FeasibleYesCount++
		case models.ResponseNoAgreement:
			nsbResponse.FeasibleNo = true
			results.Totals.FeasibleNoCount++
		case models.ResponseAbstention:
			nsbResponse.Abstention = true
			results.Totals.AbstentionCount++
		}

		// Set specific response types
		if response.Response == models.ResponseAgreeAdvance {
			nsbResponse.AcceptedAsNWIP = "Y"
			results.Totals.AcceptedAsNWIPCount++
		} else if response.Response == models.ResponseAgreeAcceptWorkingDraft {
			nsbResponse.AcceptedAsWD = "Y"
			results.Totals.AcceptedAsWDCount++
		} else if response.Response == models.ResponseAgreeCirculateCD {
			nsbResponse.AcceptedAsCD = "Y"
			results.Totals.AcceptedAsCDCount++
		} else if response.Response == models.ResponseAgreeCirculateDARF {
			nsbResponse.AcceptedAsDARF = "Y"
			results.Totals.AcceptedAsDARFCount++
		}

		// Track participation
		if response.IsCommittedToParticipate {
			results.Totals.ParticipationCount++
		}

		// Track comments
		if response.Comments != "" {
			results.Totals.CommentsCount++
		}

		results.IndividualNSBResponses = append(results.IndividualNSBResponses, nsbResponse)
	}

	// Calculate total valid responses (excluding abstentions)
	results.Totals.TotalResponses = len(*acceptance.Submissions)
	results.Totals.ValidResponses = results.Totals.TotalResponses - results.Totals.AbstentionCount

	return results, nil
}
