package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LibraryRepository struct {
	db *gorm.DB
}

func NewLibraryRepository(db *gorm.DB) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *LibraryRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *LibraryRepository) GetTopStandards(limit, offset int) ([]models.ProjectDTO, int64, error) {
	var projects []models.ProjectDTO
	var total int64

	query := r.db.Model(&models.Project{})
	query.Count(&total)

	result := query.
		Select("id, title, reference, published, created_at, updated_at").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return projects, total, nil
}

func (r *LibraryRepository) GetLatestStandards(limit, offset int) ([]models.ProjectDTO, int64, error) {
	return r.GetTopStandards(limit, offset)
}

func (r *LibraryRepository) GetTopCommittees(limit, offset int) ([]models.CommitteeDTO, int64, error) {
	var technicalCommittees []models.TechnicalCommittee
	var total int64

	// Count total committees
	if err := r.db.Model(&models.TechnicalCommittee{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch committees with preloaded Chairperson
	result := r.db.
		Preload("Chairperson.NationalStandardBody").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&technicalCommittees)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Convert to CommitteeDTO and compute counts
	committeeDTOs := make([]models.CommitteeDTO, len(technicalCommittees))
	for i, committee := range technicalCommittees {
		// Count working groups
		var workingGroupCount int64
		if err := r.db.Model(&models.WorkingGroup{}).
			Where("parent_tc_id = ?", committee.ID).
			Count(&workingGroupCount).Error; err != nil {
			return nil, 0, err
		}

		// Count current members
		var memberCount int64
		if err := r.db.Model(&models.Member{}).
			Joins("JOIN current_members cm ON cm.member_id = members.id").
			Where("cm.technical_committee_id = ?", committee.ID).
			Count(&memberCount).Error; err != nil {
			return nil, 0, err
		}

		// Count working group experts (distinct)
		var workingMemberCount int64
		if err := r.db.Model(&models.Member{}).
			Joins("JOIN working_group_experts wge ON wge.member_id = members.id").
			Joins("JOIN working_groups wg ON wge.working_group_id = wg.id").
			Where("wg.parent_tc_id = ?", committee.ID).
			Distinct("members.id").
			Count(&workingMemberCount).Error; err != nil {
			return nil, 0, err
		}

		// Count active projects
		var activeProjectCount int64
		if err := r.db.Model(&models.Project{}).
			Where("technical_committee_id = ? AND published = ? AND cancelled = ?", committee.ID, true, false).
			Count(&activeProjectCount).Error; err != nil {
			return nil, 0, err
		}

		// Map Chairperson to MemberMinified
		var chairpersonMinified *models.MemberMinified
		if committee.Chairperson != nil {
			chairpersonMinified = &models.MemberMinified{
				ID:                     committee.Chairperson.ID,
				FirstName:              committee.Chairperson.FirstName,
				LastName:               committee.Chairperson.LastName,
				NationalStandardBodyID: committee.Chairperson.NationalStandardBodyID,
				NationalStandardBody:   committee.Chairperson.NationalStandardBody,
			}
		}

		committeeDTOs[i] = models.CommitteeDTO{
			ID:                 committee.ID,
			Name:               committee.Name,
			Code:               committee.Code,
			ChairpersonId:      committee.ChairpersonId,
			Chairperson:        chairpersonMinified,
			WorkingGroupCount:  workingGroupCount,
			MemberCount:        memberCount,
			WorkingMemberCount: workingMemberCount,
			ActiveProjectCount: activeProjectCount,
		}
	}

	return committeeDTOs, total, nil
}
func (r *LibraryRepository) FindStandards(params map[string]any, limit, offset int) ([]models.ProjectDTO, int64, error) {
	var standards []models.ProjectDTO
	var total int64

	query := r.db.Model(&models.Project{}).Where("published = ?", true)

	if sector, ok := params["sector"].(string); ok && sector != "" {
		query = query.Where("sector = ?", models.ProjectSector(sector))
	}

	if title, ok := params["title"].(string); ok && title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if projectType, ok := params["type"].(models.ProjectType); ok && projectType != "" {
		query = query.Where("type = ?", projectType)
	}

	if committeeID, ok := params["committee_id"].(uuid.UUID); ok && committeeID != uuid.Nil {
		query = query.Where("technical_committee_id = ?", committeeID)
	}

	if workingGroupID, ok := params["working_group_id"].(uuid.UUID); ok && workingGroupID != uuid.Nil {
		query = query.Where("working_group_id = ?", workingGroupID)
	}

	if visible, ok := params["visible_on_library"].(bool); ok {
		query = query.Where("visible_on_library = ?", visible)
	}

	if emergency, ok := params["is_emergency"].(bool); ok {
		query = query.Where("is_emergency = ?", emergency)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := query.Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&standards)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return standards, total, nil
}

func (r *LibraryRepository) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		First(&project, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, result.Error
	}
	return &project, nil
}

func (r *LibraryRepository) GetProjectByReference(reference string) (*models.Project, error) {
	var project models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		First(&project, "reference = ?", reference)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, result.Error
	}
	return &project, nil
}

func (r *LibraryRepository) SearchProjects(query string, limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64
	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Project{}).
		Where("published = ? AND (title ILIKE ? OR description ILIKE ? OR reference ILIKE ?)", true, searchQuery, searchQuery, searchQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("published = ? AND (title ILIKE ? OR description ILIKE ? OR reference ILIKE ?)", true, searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return projects, total, nil
}

func (r *LibraryRepository) GetProjectsCreatedBetween(startDate, endDate time.Time) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("published = ? AND created_at BETWEEN ? AND ?", true, startDate, endDate).
		Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *LibraryRepository) CountProjects() (int64, error) {
	var count int64
	err := r.db.Model(&models.Project{}).Where("published = ?", true).Count(&count).Error
	return count, err
}

func (r *LibraryRepository) GetCommitteeByID(id uuid.UUID) (*models.TechnicalCommitteeDTO, error) {
	var committee models.TechnicalCommittee
	result := r.db.Preload("Chairperson").Preload("Secretary").
		Preload("WorkingGroups").Preload("SubCommittees").Preload("CurrentMembers").
		First(&committee, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("committee not found")
		}
		return nil, result.Error
	}

	// Count working groups
	var workingGroupCount int64
	if err := r.db.Model(&models.WorkingGroup{}).Where("parent_tc_id = ?", id).Count(&workingGroupCount).Error; err != nil {
		return nil, err
	}

	// Count active projects
	var activeProjectCount int64
	if err := r.db.Model(&models.Project{}).
		Where("technical_committee_id = ? AND published = ? AND cancelled = ?", id, true, false).
		Count(&activeProjectCount).Error; err != nil {
		return nil, err
	}

	// Map Chairperson to MemberMinified
	var chairpersonMinified *models.MemberMinified
	if committee.Chairperson != nil {
		chairpersonMinified = &models.MemberMinified{
			ID:                     committee.Chairperson.ID,
			FirstName:              committee.Chairperson.FirstName,
			LastName:               committee.Chairperson.LastName,
			NationalStandardBodyID: committee.Chairperson.NationalStandardBodyID,
			NationalStandardBody:   committee.Chairperson.NationalStandardBody,
		}
	}
	// Map CurrentMembers to MemberMinified
	currentMembersMinified := make([]*models.MemberMinified, len(committee.CurrentMembers))
	for i, member := range committee.CurrentMembers {
		currentMembersMinified[i] = &models.MemberMinified{
			ID:                     member.ID,
			FirstName:              member.FirstName,
			LastName:               member.LastName,
			NationalStandardBodyID: member.NationalStandardBodyID,
			NationalStandardBody:   member.NationalStandardBody,
		}
	}

	// Map to TechnicalCommitteeDTO
	committeeDTO := &models.TechnicalCommitteeDTO{
		CommitteeDTO: models.CommitteeDTO{
			ID:                 committee.ID,
			Name:               committee.Name,
			Code:               committee.Code,
			Chairperson:        chairpersonMinified,
			ChairpersonId:      committee.ChairpersonId,
			WorkingGroupCount:  workingGroupCount,
			ActiveProjectCount: activeProjectCount,
		},
		Scope:       committee.Scope,
		WorkProgram: committee.WorkProgram,
	}

	return committeeDTO, nil
}

func (r *LibraryRepository) GetCommitteeByCode(code string) (*models.TechnicalCommitteeDTO, error) {
	var committee models.TechnicalCommittee
	result := r.db.Preload("Chairperson").Preload("Secretary").
		Preload("WorkingGroups").Preload("SubCommittees").Preload("CurrentMembers").
		First(&committee, "code = ?", code)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("committee not found")
		}
		return nil, result.Error
	}

	// Count working groups
	var workingGroupCount int64
	if err := r.db.Model(&models.WorkingGroup{}).Where("parent_tc_id = ?", committee.ID).Count(&workingGroupCount).Error; err != nil {
		return nil, err
	}

	// Count active projects
	var activeProjectCount int64
	if err := r.db.Model(&models.Project{}).
		Where("technical_committee_id = ? AND published = ? AND cancelled = ?", committee.ID, true, false).
		Count(&activeProjectCount).Error; err != nil {
		return nil, err
	}

	// Map Chairperson to MemberMinified
	var chairpersonMinified *models.MemberMinified
	if committee.Chairperson != nil {
		chairpersonMinified = &models.MemberMinified{
			ID:                     committee.Chairperson.ID,
			FirstName:              committee.Chairperson.FirstName,
			LastName:               committee.Chairperson.LastName,
			NationalStandardBodyID: committee.Chairperson.NationalStandardBodyID,
			NationalStandardBody:   committee.Chairperson.NationalStandardBody,
		}
	}

	// Map CurrentMembers to MemberMinified
	currentMembersMinified := make([]*models.MemberMinified, len(committee.CurrentMembers))
	for i, member := range committee.CurrentMembers {
		currentMembersMinified[i] = &models.MemberMinified{
			ID:                     member.ID,
			FirstName:              member.FirstName,
			LastName:               member.LastName,
			NationalStandardBodyID: member.NationalStandardBodyID,
			NationalStandardBody:   member.NationalStandardBody,
		}
	}

	// Map to TechnicalCommitteeDTO
	committeeDTO := &models.TechnicalCommitteeDTO{
		CommitteeDTO: models.CommitteeDTO{
			ID:                 committee.ID,
			Name:               committee.Name,
			Code:               committee.Code,
			Chairperson:        chairpersonMinified,
			ChairpersonId:      committee.ChairpersonId,
			WorkingGroupCount:  workingGroupCount,
			ActiveProjectCount: activeProjectCount,
		},
		Scope:       committee.Scope,
		WorkProgram: committee.WorkProgram,
	}

	return committeeDTO, nil
}

func (r *LibraryRepository) ListCommittees(limit, offset int) ([]models.TechnicalCommitteeDTO, int64, error) {
	var committees []models.TechnicalCommittee
	var total int64

	err := r.db.Model(&models.TechnicalCommittee{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Chairperson").Preload("Secretary").
		Preload("WorkingGroups").Preload("SubCommittees").Preload("CurrentMembers").
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&committees)
	if result.Error != nil {
		return nil, 0, err
	}

	// Convert to TechnicalCommitteeDTO and fetch counts
	committeeDTOs := make([]models.TechnicalCommitteeDTO, len(committees))
	for i, committee := range committees {
		// Count working groups
		var workingGroupCount int64
		if err := r.db.Model(&models.WorkingGroup{}).Where("parent_tc_id = ?", committee.ID).Count(&workingGroupCount).Error; err != nil {
			return nil, 0, err
		}

		// Count active projects
		var activeProjectCount int64
		if err := r.db.Model(&models.Project{}).
			Where("technical_committee_id = ? AND published = ? AND cancelled = ?", committee.ID, true, false).
			Count(&activeProjectCount).Error; err != nil {
			return nil, 0, err
		}

		// Map Chairperson to MemberMinified
		var chairpersonMinified *models.MemberMinified
		if committee.Chairperson != nil {
			chairpersonMinified = &models.MemberMinified{
				ID:                     committee.Chairperson.ID,
				FirstName:              committee.Chairperson.FirstName,
				LastName:               committee.Chairperson.LastName,
				NationalStandardBodyID: committee.Chairperson.NationalStandardBodyID,
				NationalStandardBody:   committee.Chairperson.NationalStandardBody,
			}
		}
		// Map CurrentMembers to MemberMinified
		currentMembersMinified := make([]*models.MemberMinified, len(committee.CurrentMembers))
		for j, member := range committee.CurrentMembers {
			currentMembersMinified[j] = &models.MemberMinified{
				ID:                     member.ID,
				FirstName:              member.FirstName,
				LastName:               member.LastName,
				NationalStandardBodyID: member.NationalStandardBodyID,
				NationalStandardBody:   member.NationalStandardBody,
			}
		}

		committeeDTOs[i] = models.TechnicalCommitteeDTO{
			CommitteeDTO: models.CommitteeDTO{
				ID:                 committee.ID,
				Name:               committee.Name,
				Code:               committee.Code,
				Chairperson:        chairpersonMinified,
				ChairpersonId:      committee.ChairpersonId,
				WorkingGroupCount:  workingGroupCount,
				ActiveProjectCount: activeProjectCount,
			},
			Scope:       committee.Scope,
			WorkProgram: committee.WorkProgram,
		}
	}

	return committeeDTOs, total, nil
}

func (r *LibraryRepository) SearchCommittees(query string, limit, offset int) ([]models.TechnicalCommitteeDTO, int64, error) {
	var committees []models.TechnicalCommittee
	var total int64
	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.TechnicalCommittee{}).
		Where("name ILIKE ? OR code ILIKE ?", searchQuery, searchQuery)
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Chairperson").Preload("Secretary").
		Preload("WorkingGroups").Preload("SubCommittees").Preload("CurrentMembers").
		Where("name ILIKE ? OR code ILIKE ?", searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&committees)
	if result.Error != nil {
		return nil, 0, err
	}

	// Convert to TechnicalCommitteeDTO and fetch counts
	committeeDTOs := make([]models.TechnicalCommitteeDTO, len(committees))
	for i, committee := range committees {
		// Count working groups
		var workingGroupCount int64
		if err := r.db.Model(&models.WorkingGroup{}).Where("parent_tc_id = ?", committee.ID).Count(&workingGroupCount).Error; err != nil {
			return nil, 0, err
		}

		// Count active projects
		var activeProjectCount int64
		if err := r.db.Model(&models.Project{}).
			Where("technical_committee_id = ? AND published = ? AND cancelled = ?", committee.ID, true, false).
			Count(&activeProjectCount).Error; err != nil {
			return nil, 0, err
		}

		// Map Chairperson to MemberMinified
		var chairpersonMinified *models.MemberMinified
		if committee.Chairperson != nil {
			chairpersonMinified = &models.MemberMinified{
				ID:                     committee.Chairperson.ID,
				FirstName:              committee.Chairperson.FirstName,
				LastName:               committee.Chairperson.LastName,
				NationalStandardBodyID: committee.Chairperson.NationalStandardBodyID,
				NationalStandardBody:   committee.Chairperson.NationalStandardBody,
			}
		}
		// Map CurrentMembers to MemberMinified
		currentMembersMinified := make([]*models.MemberMinified, len(committee.CurrentMembers))
		for j, member := range committee.CurrentMembers {
			currentMembersMinified[j] = &models.MemberMinified{
				ID:                     member.ID,
				FirstName:              member.FirstName,
				LastName:               member.LastName,
				NationalStandardBodyID: member.NationalStandardBodyID,
				NationalStandardBody:   member.NationalStandardBody,
			}
		}

		committeeDTOs[i] = models.TechnicalCommitteeDTO{
			CommitteeDTO: models.CommitteeDTO{
				ID:                 committee.ID,
				Name:               committee.Name,
				Code:               committee.Code,
				Chairperson:        chairpersonMinified,
				ChairpersonId:      committee.ChairpersonId,
				WorkingGroupCount:  workingGroupCount,
				ActiveProjectCount: activeProjectCount,
			},
			Scope:       committee.Scope,
			WorkProgram: committee.WorkProgram,
		}
	}

	return committeeDTOs, total, nil
}

func (r *LibraryRepository) CountCommittees() (int64, error) {
	var count int64
	err := r.db.Model(&models.TechnicalCommittee{}).Count(&count).Error
	return count, err
}

func (r *LibraryRepository) GetProjectsByCommitteeID(committeeID string) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("published = ? AND technical_committee_id = ?", true, committeeID).
		Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}
func (r *LibraryRepository) GetSectors() ([]models.ProjectSector, error) {
	return []models.ProjectSector{
		models.Health,
		models.IT,
		models.Management,
		models.Safety,
		models.Transport,
		models.Energy,
		models.Diversity,
		models.Environment,
		models.Food,
		models.Building,
		models.Engineering,
	}, nil
}
