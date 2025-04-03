package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentRepository handles database operations for Document entities
type DocumentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository creates a new document repository instance
func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create adds a new document to the database
func (r *DocumentRepository) Create(doc *models.Document) error {
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}
	return r.db.Create(doc).Error
}

func (r *DocumentRepository) UpdateProjectDoc(projectId, docType, fileURL, member string) error {
	tx := r.db.Begin() // Start a transaction
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var project models.Project
	if err := tx.Where("id = ?", projectId).First(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	doc := models.Document{
		ID:          uuid.New(),
		Reference:   project.Reference,
		Title:       project.Reference,
		FileURL:     fileURL,
		Description: docType,
		CreatedByID: member,
	}

	if err := tx.Create(&doc).Error; err != nil {
		tx.Rollback()
		return err
	}

	if docType == "WD" {
		docID := doc.ID.String()
		project.WorkingDraftID = &docID
	}
	if docType == "CD" {
		docID := doc.ID.String()
		project.CommitteeDraftID = &docID
	}

	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetByID retrieves a document by its ID
func (r *DocumentRepository) GetByID(id uuid.UUID) (*models.Document, error) {
	var doc models.Document
	err := r.db.First(&doc, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return &doc, nil
}

// GetByReference retrieves a document by its reference
func (r *DocumentRepository) GetByReference(reference string) (*models.Document, error) {
	var doc models.Document
	err := r.db.First(&doc, "reference = ?", reference).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return &doc, nil
}

// GetByTitle retrieves a document by its title
func (r *DocumentRepository) GetByTitle(title string) (*models.Document, error) {
	var doc models.Document
	err := r.db.First(&doc, "title = ?", title).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return &doc, nil
}

// Update updates an existing document
func (r *DocumentRepository) Update(doc *models.Document) error {
	return r.db.Save(doc).Error
}

// UpdatePartial updates specific fields of a document
func (r *DocumentRepository) UpdatePartial(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&models.Document{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateFileURL updates only the file URL of a document
func (r *DocumentRepository) UpdateFileURL(id uuid.UUID, fileURL string) error {
	return r.db.Model(&models.Document{}).Where("id = ?", id).Update("file_url", fileURL).Error
}

// Delete removes a document from the database
func (r *DocumentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Document{}, "id = ?", id).Error
}

// List retrieves all documents with pagination
func (r *DocumentRepository) List(limit, offset int) ([]models.Document, int64, error) {
	var docs []models.Document
	var total int64

	// Get total count
	if err := r.db.Model(&models.Document{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get documents with pagination
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&docs).Error
	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// Search searches for documents by title, description or reference
func (r *DocumentRepository) Search(query string, limit, offset int) ([]models.Document, int64, error) {
	var docs []models.Document
	var total int64

	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Document{}).
		Where("title LIKE ? OR description LIKE ? OR reference LIKE ?",
			searchQuery, searchQuery, searchQuery)

	// Get total count
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get documents with pagination
	err := r.db.Where("title LIKE ? OR description LIKE ? OR reference LIKE ?",
		searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&docs).Error

	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// GetDocumentsCreatedBetween retrieves documents created within a time range
func (r *DocumentRepository) GetDocumentsCreatedBetween(startDate, endDate time.Time) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&docs).Error
	return docs, err
}

// Exists checks if a document exists by id or reference or title
func (r *DocumentRepository) Exists(id uuid.UUID, reference string, title string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Document{}).
		Where("id = ? OR reference = ? OR title = ?", id, reference, title).
		Count(&count).Error
	return count > 0, err
}

// CountAll returns the total number of documents
func (r *DocumentRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Document{}).Count(&count).Error
	return count, err
}

func (r *DocumentRepository) ProjectDocuments(projectId string) ([]models.Document, error) {
	var project models.Project
	var docs []models.Document
	if err := r.db.Where("id = ?", projectId).Preload("WorkingDraft").Preload("CommitteeDraft").First(&project).Error; err != nil {
		return docs, err
	}
	if project.WorkingDraftID != nil {
		docs = append(docs, *project.WorkingDraft)
	}
	if project.CommitteeDraftID != nil {
		docs = append(docs, *project.CommitteeDraft)
	}
	return docs, nil
}
