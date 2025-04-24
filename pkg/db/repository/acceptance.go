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

		var member models.Member
		if err := tx.Where("id = ?", response.ResponderID).Preload(clause.Associations).First(&member).Error; err != nil {
			tx.Rollback()
			return err
		}

		// if member.NationalStandardBody != nil && member.NationalStandardBody.NationalTCSecretaryID != &response.ResponderID {
		// 	tx.Rollback()
		// 	return errors.New("user is not a National TC secretary of the responding NSB")
		// }

		// Check if Acceptance exists for the given project
		if err := tx.Where("project_id = ?", response.ProjectID).First(&acceptance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Get Project
				var project models.Project
				if err := tx.Where("id = ?", response.ProjectID).Preload("TechnicalCommittee").First(&project).Error; err != nil {
					tx.Rollback()
					return err
				}

				// Create a new Acceptance if it does not exist
				acceptance = models.Acceptance{
					ID:            uuid.New(),
					ProjectID:     response.ProjectID,
					IsApproved:    false,
					CreatedAt:     time.Now(),
					TCSecretaryID: project.TechnicalCommittee.SecretaryId,
				}
				if err := tx.Create(&acceptance).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return err
			}
		}

		// Attach NSBResponse to the existing acceptance
		response.AcceptanceID = acceptance.ID
		nsb := member.NationalStandardBodyID
		response.RespondingNSBID = *nsb
		response.NationalTCSecretaryID = member.NationalStandardBody.NationalTCSecretaryID

		// Process relevant standards from string IDs to Document associations
		if len(response.RelevantStandards) > 0 {
			var standardDocs []models.Document
			if err := tx.Where("id IN ?", response.RelevantStandards).Find(&standardDocs).Error; err != nil {
				tx.Rollback()
				return err
			}
			response.RelevantStandardsRefs = &standardDocs
		}

		// Process relevant regulations from string IDs to Document associations
		if len(response.RelevantRegulations) > 0 {
			var regulationDocs []models.Document
			if err := tx.Where("id IN ?", response.RelevantRegulations).Find(&regulationDocs).Error; err != nil {
				tx.Rollback()
				return err
			}
			response.RelevantRegulationsRefs = &regulationDocs
		}

		if err := tx.Create(response).Error; err != nil {
			tx.Rollback()
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
	if err := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Preload(clause.Associations).Preload("RespondingNSB.MemberState").Find(&NSBResponses).Error; err != nil {
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
	if err := r.db.Where("project_id = ?", id).First(&acceptance).Error; err != nil {
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

func (r *AcceptanceRepository) SetAcceptanceApproval(results models.Acceptance) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Where("id = ?", results.ProjectID).Preload("TechnicalCommittee").First(&project).Error; err != nil {
			return err
		}

		if project.TechnicalCommittee.SecretaryId == nil || *project.TechnicalCommittee.SecretaryId != *results.TCSecretaryID {
			tx.Rollback()
			return fmt.Errorf("User is not allowed to perform this action")
		}
		var acceptance models.Acceptance
		if err := tx.Where("project_id = ?", results.ProjectID).First(&acceptance).Error; err != nil {
			tx.Rollback()
			return err
		}

		acceptance.IsApproved = true
		acceptance.ApprovalCriteriaMet = true
		now := time.Now()
		acceptance.SMCApprovalDate = &now

		acceptance.DevelopmentTrack = results.DevelopmentTrack
		acceptance.DraftStatus = results.DraftStatus
		acceptance.IsPreliminaryWork = results.IsPreliminaryWork
		acceptance.IsActiveWork = results.IsActiveWork
		acceptance.TCSecretaryID = results.TCSecretaryID
		acceptance.OtherInformation = results.OtherInformation
		acceptance.DraftExpectedDate = results.DraftExpectedDate

		if len(results.DocumentsInConsidaration) > 0 {
			var documents []models.Document
			if err := tx.Where("id IN (?)", results.DocumentsInConsidaration).Find(&documents).Error; err != nil {
				tx.Rollback()
				return err
			}
			acceptance.DocumentsToConsider = &documents
		}

		if err := tx.Save(&acceptance).Error; err != nil {
			tx.Rollback()
			return err
		}

		var stage models.Stage
		if err := tx.Where("number = ?", 2).First(&stage).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Update the project stage
		if err := UpdateProjectStageWithTx(tx, results.ProjectID, stage.ID.String(), "Proposal Accepted", "NWIP", stage.Abbreviation); err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func (r *AcceptanceRepository) GetAcceptanceResults(id string) (*models.AcceptanceResults, error) {
	// Get Project Proposal
	var proposal models.Proposal
	if err := r.db.Where("project_id = ?", id).First(&proposal).Error; err != nil {
		return nil, err
	}
	// Find the Acceptance by Project ID
	var acceptance models.Acceptance
	if err := r.db.Where("project_id = ?", id).Preload("Submissions").Preload("Submissions.RespondingNSB").First(&acceptance).Error; err != nil {
		return nil, err
	}

	// Initialize results structure
	results := &models.AcceptanceResults{
		AcceptanceID:           id,
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

	criteriaMet, comment, err := checkIfCriteriaIsMet(proposal.ExistingIntlStandard, results.IndividualNSBResponses)
	if err != nil {
		return nil, err
	}

	// Check if criteria is met
	results.CriterialMet = criteriaMet
	results.CriterialComments = comment

	return results, nil
}

func checkIfCriteriaIsMet(internationalStandard bool, responses []models.IndividualNSBResponse) (bool, string, error) {
	// Count total P-members voting (excluding abstentions)
	var totalVotingCount int
	var favorableVoteCount int
	var favorableParticipationCount int // Participants who voted in favor
	var generalParticipationCount int   // All participants regardless of vote

	for _, response := range responses {
		// Skip abstentions when counting total voting members
		if response.Abstention {
			continue
		}

		totalVotingCount++

		// Count members willing to participate actively
		if response.Participation {
			generalParticipationCount++

			// Count participating members who also voted favorably
			if response.FeasibleYes {
				favorableParticipationCount++
			}
		}

		// Count favorable votes (FeasibleYes)
		if response.FeasibleYes {
			favorableVoteCount++
		}
	}

	// If there are no votes, return false
	if totalVotingCount < 6 {
		return false, "At least 6 votes are required", nil
	}

	if internationalStandard {
		// For international standards: MORE than 50% of P-members voting in favor
		percentageInFavor := float64(favorableVoteCount) / float64(totalVotingCount) * 100

		// Check if more than 50% voted in favor
		if percentageInFavor <= 50.0 {
			return false, fmt.Sprintf("Criteria not met: Only %.1f%% of P-members voted in favor (more than 50%% required)", percentageInFavor), nil
		}

		// Verify there are participants who voted in favor
		if favorableParticipationCount == 0 {
			return false, "Criteria not met: No members who voted in favor committed to participate", nil
		}

		return true, fmt.Sprintf("Criteria met: %.1f%% of P-members voted in favor with %d favorable members committed to participate", percentageInFavor, favorableParticipationCount), nil
	} else {
		// For projects requiring preparatory/committee stages:
		// Need at least 6 P-members voting and at least 3 members willing to participate actively

		if totalVotingCount < 6 {
			return false, fmt.Sprintf("Criteria not met: Only %d P-members voted (minimum 6 required)", totalVotingCount), nil
		}

		if generalParticipationCount < 3 {
			return false, fmt.Sprintf("Criteria not met: Only %d members willing to participate actively (minimum 3 required)", generalParticipationCount), nil
		}

		return true, fmt.Sprintf("Criteria met: %d P-members voted and %d members willing to participate actively", totalVotingCount, generalParticipationCount), nil
	}
}
