package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NSBResponseStatusChangeRepository struct {
	db *gorm.DB
}

func NewNSBResponseStatusChangeRepository(db *gorm.DB) *NSBResponseStatusChangeRepository {
	return &NSBResponseStatusChangeRepository{
		db: db,
	}
}

func (r *NSBResponseStatusChangeRepository) Create(change *models.NSBResponseStatusChange) error {
	if change.ID == uuid.Nil {
		change.ID = uuid.New()
	}
	change.CreatedAt = time.Now()

	// Default status if not set
	if change.Status == "" {
		change.Status = models.PENDING
	}

	return r.db.Create(change).Error
}

func (r *NSBResponseStatusChangeRepository) GetByID(id uuid.UUID) (*models.NSBResponseStatusChange, error) {
	var change models.NSBResponseStatusChange

	result := r.db.
		Preload("Responder").
		Preload("InitialResponse").
		Preload("TCSecretariat").
		Where("id = ?", id).
		First(&change)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("NSBResponseStatusChange not found")
		}
		return nil, result.Error
	}

	return &change, nil
}

func (r *NSBResponseStatusChangeRepository) GetByInitialResponseID(initialResponseID string) ([]models.NSBResponseStatusChange, error) {
	var changes []models.NSBResponseStatusChange

	result := r.db.
		Preload("Responder").
		Preload("InitialResponse").
		Preload("TCSecretariat").
		Where("initial_response_id = ?", initialResponseID).
		Order("created_at DESC").
		Find(&changes)

	if result.Error != nil {
		return nil, result.Error
	}

	return changes, nil
}

func (r *NSBResponseStatusChangeRepository) GetByResponderID(responderID string) ([]models.NSBResponseStatusChange, error) {
	var changes []models.NSBResponseStatusChange

	result := r.db.
		Preload("Responder").
		Preload("InitialResponse").
		Preload("TCSecretariat").
		Where("responder_id = ?", responderID).
		Order("created_at DESC").
		Find(&changes)

	if result.Error != nil {
		return nil, result.Error
	}

	return changes, nil
}

func (r *NSBResponseStatusChangeRepository) GetPendingChanges() ([]models.NSBResponseStatusChange, error) {
	var changes []models.NSBResponseStatusChange

	result := r.db.
		Preload("Responder").
		Preload("InitialResponse").
		Preload("TCSecretariat").
		Where("status = ?", models.PENDING).
		Order("created_at ASC").
		Find(&changes)

	if result.Error != nil {
		return nil, result.Error
	}

	return changes, nil
}

func (r *NSBResponseStatusChangeRepository) UpdateStatus(id uuid.UUID, status models.Status, tcSecretariatID string) error {
	return r.db.Model(&models.NSBResponseStatusChange{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            status,
			"tc_secretariat_id": tcSecretariatID,
		}).Error
}

func (r *NSBResponseStatusChangeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.NSBResponseStatusChange{}, id).Error
}

func (r *NSBResponseStatusChangeRepository) GetByStatus(status models.Status) ([]models.NSBResponseStatusChange, error) {
	var changes []models.NSBResponseStatusChange

	result := r.db.
		Preload("Responder").
		Preload("InitialResponse").
		Preload("TCSecretariat").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&changes)

	if result.Error != nil {
		return nil, result.Error
	}

	return changes, nil
}

func (r *NSBResponseStatusChangeRepository) ApproveChange(id uuid.UUID, tcSecretariatID string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get the change request with initial response info
	var change models.NSBResponseStatusChange
	if err := tx.Preload("InitialResponse").Where("id = ?", id).First(&change).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the change request status
	if err := tx.Model(&change).Updates(map[string]interface{}{
		"status":            models.APPROVED,
		"tc_secretariat_id": tcSecretariatID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the initial NSBResponse with the new values
	if err := tx.Model(&models.NSBResponse{}).Where("id = ?", change.InitialResponseID).Updates(map[string]interface{}{
		"response":                    change.Response,
		"is_committed_to_participate": change.IsCommittedToParticipate,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *NSBResponseStatusChangeRepository) RejectChange(id uuid.UUID, tcSecretariatID, comment string) error {
	return r.db.Model(&models.NSBResponseStatusChange{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":                 models.REJECTED,
			"tc_secretariat_id":      tcSecretariatID,
			"tc_secretariat_comment": comment,
		}).Error
}
