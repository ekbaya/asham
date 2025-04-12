package repository

import (
	"errors"
	"fmt"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) CreateMember(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *MemberRepository) GetMemberByID(id string) (*models.Member, error) {
	var member models.Member
	result := r.db.Preload("NationalStandardBody").First(&member, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func (r *MemberRepository) GetAllMembers() (*[]models.Member, error) {
	var members []models.Member
	result := r.db.Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return &members, nil
}

func (r *MemberRepository) GetMemberByEmail(email string) (*models.Member, error) {
	var member models.Member
	result := r.db.Where("email = ?", email).First(&member)

	if result.Error != nil {
		// Check if the error is "record not found"
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of an error when member not found
		}
		return nil, result.Error // Return other errors as is
	}

	return &member, nil
}

func (r *MemberRepository) EmailExists(email string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Member{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *MemberRepository) GetMembersByCountryCode(countryCode string) (*[]models.Member, error) {
	var members []models.Member
	result := r.db.Where("country_code = ?", countryCode).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return &members, nil
}

func (r *MemberRepository) UpdateMember(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *MemberRepository) DeleteMember(memberID string) error {
	return r.db.Delete(&models.Member{}, "id = ?", memberID).Error
}

func (r *MemberRepository) GetMemberResponsibilities(memberID string) (map[string]any, error) {
	// Initialize the result map
	responsibilities := make(map[string]any)

	// Get the member to ensure it exists
	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	// Technical Committees where member is present
	var technicalCommittees []models.TechnicalCommittee
	if err := r.db.Joins("JOIN current_members ON current_members.technical_committee_id = technical_committees.id").
		Where("current_members.member_id = ?", memberID).
		Find(&technicalCommittees).Error; err != nil {
		return nil, err
	}

	tcList := make([]map[string]any, 0)
	for _, tc := range technicalCommittees {
		role := "Member"
		if tc.ChairpersonId != nil && *tc.ChairpersonId == memberID {
			role = "Chairperson"
		} else if tc.SecretaryId != nil && *tc.SecretaryId == memberID {
			role = "Secretary"
		}

		tcList = append(tcList, map[string]any{
			"id":   tc.ID,
			"name": tc.Name,
			"role": role,
		})
	}
	responsibilities["technical_committees"] = tcList

	// Sub-Committees where member is present
	var subCommittees []models.SubCommittee
	if err := r.db.Joins("JOIN sc_members ON sc_members.sub_committee_id = sub_committees.id").
		Where("sc_members.member_id = ?", memberID).
		Find(&subCommittees).Error; err != nil {
		return nil, err
	}

	scList := make([]map[string]any, 0)
	for _, sc := range subCommittees {
		role := "Member"
		if sc.ChairpersonId != nil && *sc.ChairpersonId == memberID {
			role = "Chairperson"
		} else if sc.SecretaryId != nil && *sc.SecretaryId == memberID {
			role = "Secretary"
		}

		scList = append(scList, map[string]any{
			"id":        sc.ID,
			"name":      sc.Name,
			"parent_tc": sc.ParentTCID,
			"role":      role,
		})
	}
	responsibilities["sub_committees"] = scList

	// Working Groups where member is present (either as convenor or expert)
	var workingGroups []models.WorkingGroup
	if err := r.db.Joins("JOIN working_group_experts ON working_group_experts.working_group_id = working_groups.id").
		Where("working_group_experts.member_id = ? OR working_groups.convenor_id = ?", memberID, memberID).
		Find(&workingGroups).Error; err != nil {
		return nil, err
	}

	wgList := make([]map[string]any, 0)
	for _, wg := range workingGroups {
		role := "Expert"
		if wg.ConvenorId == memberID {
			role = "Convenor"
		}

		wgList = append(wgList, map[string]any{
			"id":        wg.ID,
			"name":      wg.Name,
			"parent_tc": wg.ParentTCID,
			"role":      role,
		})
	}
	responsibilities["working_groups"] = wgList

	// Task Forces where member is present (either as convenor or national delegation)
	var taskForces []models.TaskForce
	if err := r.db.Joins("JOIN national_deligations ON national_deligations.task_force_id = task_forces.id").
		Where("national_deligations.member_id = ? OR task_forces.convenor_id = ?", memberID, memberID).
		Find(&taskForces).Error; err != nil {
		return nil, err
	}

	tfList := make([]map[string]any, 0)
	for _, tf := range taskForces {
		role := "National Delegate"
		if tf.ConvenorId == memberID {
			role = "Convenor"
		}

		tfList = append(tfList, map[string]any{
			"id":        tf.ID,
			"name":      tf.Name,
			"parent_tc": tf.ParentTCID,
			"role":      role,
		})
	}
	responsibilities["task_forces"] = tfList

	// Specialized Committees
	var specializedCommittees []models.SpecializedCommittee
	if err := r.db.Joins("JOIN specialized_committee_members ON specialized_committee_members.specialized_committee_id = specialized_committees.id").
		Where("specialized_committee_members.member_id = ?", memberID).
		Find(&specializedCommittees).Error; err != nil {
		return nil, err
	}

	specList := make([]map[string]any, 0)
	for _, spec := range specializedCommittees {
		role := "Member"
		if spec.ChairpersonId != nil && *spec.ChairpersonId == memberID {
			role = "Chairperson"
		} else if spec.SecretaryId != nil && *spec.SecretaryId == memberID {
			role = "Secretary"
		}

		specList = append(specList, map[string]any{
			"id":   spec.ID,
			"name": spec.Name,
			"type": spec.Type,
			"role": role,
		})
	}
	responsibilities["specialized_committees"] = specList

	// Standards Management Committee
	var smcs []models.StandardsManagementCommittee
	if err := r.db.Model(&models.StandardsManagementCommittee{}).
		Joins("LEFT JOIN regional_representatives ON regional_representatives.standards_management_committee_id = standards_management_committees.id").
		Joins("LEFT JOIN elected_members ON elected_members.standards_management_committee_id = standards_management_committees.id").
		Joins("LEFT JOIN observers ON observers.standards_management_committee_id = standards_management_committees.id").
		Where("regional_representatives.member_id = ? OR elected_members.member_id = ? OR observers.member_id = ? OR standards_management_committees.chairperson_id = ? OR standards_management_committees.secretary_id = ?",
			memberID, memberID, memberID, memberID, memberID).
		Find(&smcs).Error; err != nil {
		return nil, err
	}

	smcList := make([]map[string]any, 0)
	for _, smc := range smcs {
		role := "Unknown"
		// Determine role
		if smc.ChairpersonId != nil && *smc.ChairpersonId == memberID {
			role = "Chairperson"
		} else if smc.SecretaryId != nil && *smc.SecretaryId == memberID {
			role = "Secretary"
		} else {
			// Check in which category the member belongs
			var regionCount int64
			r.db.Model(&models.Member{}).
				Joins("JOIN regional_representatives ON regional_representatives.member_id = members.id").
				Where("members.id = ? AND regional_representatives.standards_management_committee_id = ?", memberID, smc.ID).
				Count(&regionCount)

			if regionCount > 0 {
				role = "Regional Representative"
			} else {
				var electedCount int64
				r.db.Model(&models.Member{}).
					Joins("JOIN elected_members ON elected_members.member_id = members.id").
					Where("members.id = ? AND elected_members.standards_management_committee_id = ?", memberID, smc.ID).
					Count(&electedCount)

				if electedCount > 0 {
					role = "Elected Member"
				} else {
					var observerCount int64
					r.db.Model(&models.Member{}).
						Joins("JOIN observers ON observers.member_id = members.id").
						Where("members.id = ? AND observers.standards_management_committee_id = ?", memberID, smc.ID).
						Count(&observerCount)

					if observerCount > 0 {
						role = "Observer"
					}
				}
			}
		}

		smcList = append(smcList, map[string]any{
			"id":   smc.ID,
			"name": smc.Name,
			"role": role,
		})
	}
	responsibilities["standards_management_committees"] = smcList

	// Similar patterns for other committee types:
	// ARSO Council
	var arsoCouncils []models.ARSOCouncil
	if err := r.db.Joins("JOIN arsocouncil_members ON arsocouncil_members.arsocouncil_id = arsocouncils.id").
		Where("arsocouncil_members.member_id = ? OR arsocouncils.chairperson_id = ? OR arsocouncils.secretary_id = ?",
			memberID, memberID, memberID).
		Find(&arsoCouncils).Error; err != nil {
		return nil, err
	}

	arsoList := make([]map[string]any, 0)
	for _, arso := range arsoCouncils {
		role := "Member"
		if arso.ChairpersonId != nil && *arso.ChairpersonId == memberID {
			role = "Chairperson"
		} else if arso.SecretaryId != nil && *arso.SecretaryId == memberID {
			role = "Secretary"
		}

		arsoList = append(arsoList, map[string]any{
			"id":   arso.ID,
			"name": arso.Name,
			"role": role,
		})
	}
	responsibilities["arso_councils"] = arsoList

	// Joint Advisory Group
	var jags []models.JointAdvisoryGroup
	if err := r.db.Model(&models.JointAdvisoryGroup{}).
		Joins("LEFT JOIN jag_members ON jag_members.joint_advisory_group_id = joint_advisory_groups.id").
		Joins("LEFT JOIN jag_observers ON jag_observers.joint_advisory_group_id = joint_advisory_groups.id").
		Where("jag_members.member_id = ? OR jag_observers.member_id = ? OR joint_advisory_groups.chairperson_id = ? OR joint_advisory_groups.secretary_id = ?",
			memberID, memberID, memberID, memberID).
		Find(&jags).Error; err != nil {
		return nil, err
	}

	jagList := make([]map[string]any, 0)
	for _, jag := range jags {
		role := "Unknown"
		if jag.ChairpersonId != nil && *jag.ChairpersonId == memberID {
			role = "Chairperson"
		} else if jag.SecretaryId != nil && *jag.SecretaryId == memberID {
			role = "Secretary"
		} else {
			// Check if member or observer
			var memberCount int64
			r.db.Model(&models.Member{}).
				Joins("JOIN jag_members ON jag_members.member_id = members.id").
				Where("members.id = ? AND jag_members.joint_advisory_group_id = ?", memberID, jag.ID).
				Count(&memberCount)

			if memberCount > 0 {
				role = "Regional Economic Community Member"
			} else {
				var observerCount int64
				r.db.Model(&models.Member{}).
					Joins("JOIN jag_observers ON jag_observers.member_id = members.id").
					Where("members.id = ? AND jag_observers.joint_advisory_group_id = ?", memberID, jag.ID).
					Count(&observerCount)

				if observerCount > 0 {
					role = "Observer Member"
				}
			}
		}

		jagList = append(jagList, map[string]any{
			"id":   jag.ID,
			"name": jag.Name,
			"role": role,
		})
	}
	responsibilities["joint_advisory_groups"] = jagList

	// Joint Technical Committee
	var jtcs []models.JointTechnicalCommittee
	if err := r.db.Joins("JOIN joint_members ON joint_members.joint_technical_committee_id = joint_technical_committees.id").
		Where("joint_members.member_id = ? OR joint_technical_committees.chairperson_id = ? OR joint_technical_committees.secretary_id = ?",
			memberID, memberID, memberID).
		Find(&jtcs).Error; err != nil {
		return nil, err
	}

	jtcList := make([]map[string]any, 0)
	for _, jtc := range jtcs {
		role := "Joint Member"
		if jtc.ChairpersonId != nil && *jtc.ChairpersonId == memberID {
			role = "Chairperson"
		} else if jtc.SecretaryId != nil && *jtc.SecretaryId == memberID {
			role = "Secretary"
		}

		jtcList = append(jtcList, map[string]any{
			"id":   jtc.ID,
			"name": jtc.Name,
			"role": role,
		})
	}
	responsibilities["joint_technical_committees"] = jtcList

	return responsibilities, nil
}
