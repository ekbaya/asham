package services

import (
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type DocumentService struct {
	repo        *repository.DocumentRepository
	projectRepo *repository.ProjectRepository
}

func NewDocumentService(repo *repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (service *DocumentService) Create(doc *models.Document) error {
	doc.ID = uuid.New()
	doc.CreatedAt = time.Now()
	return service.repo.Create(doc)
}

func (service *DocumentService) UpdateProjectDoc(project, docType, fileURL, member string) error {
	return service.repo.UpdateProjectDoc(project, docType, fileURL, member)
}

func (service *DocumentService) GetByID(id uuid.UUID) (*models.Document, error) {
	return service.repo.GetByID(id)
}

func (service *DocumentService) GetByReference(reference string) (*models.Document, error) {
	return service.repo.GetByReference(reference)
}

func (service *DocumentService) GetByTitle(title string) (*models.Document, error) {
	return service.repo.GetByTitle(title)
}

func (service *DocumentService) Update(doc *models.Document) error {
	return service.repo.Update(doc)
}

func (service *DocumentService) UpdatePartial(id uuid.UUID, updates map[string]interface{}) error {
	return service.repo.UpdatePartial(id, updates)
}

func (service *DocumentService) UpdateFileURL(id uuid.UUID, fileURL string) error {
	return service.repo.UpdateFileURL(id, fileURL)
}

func (service *DocumentService) Delete(id uuid.UUID) error {
	// Check if document is attached to a project e.g wd/cd and set to nil
	projects, err := service.projectRepo.FindByDocumentID(id)
	if err != nil {
		return fmt.Errorf("error finding projects with document %s: %w", id, err)
	}

	docId := id.String()

	// Update any projects that reference this document
	for _, project := range projects {
		updated := false

		if project.WorkingDraftID == &docId {
			project.WorkingDraftID = nil
			updated = true
		}

		if project.CommitteeDraftID == &docId {
			project.CommitteeDraftID = nil
			updated = true
		}

		if updated {
			if err := service.projectRepo.Update(&project); err != nil {
				return fmt.Errorf("error updating project %s: %w", project.ID, err)
			}
		}
	}
	// Now delete the document
	return service.repo.Delete(id)
}

func (service *DocumentService) Exists(id uuid.UUID, reference string, title string) (bool, error) {
	return service.repo.Exists(id, reference, title)
}

func (service *DocumentService) List(limit, offset int) ([]models.Document, int64, error) {
	return service.repo.List(limit, offset)
}

func (service *DocumentService) Search(query string, limit, offset int) ([]models.Document, int64, error) {
	return service.repo.Search(query, limit, offset)
}

func (service *DocumentService) GetDocumentsCreatedBetween(startDate, endDate time.Time) ([]models.Document, error) {
	return service.repo.GetDocumentsCreatedBetween(startDate, endDate)
}

func (service *DocumentService) CountAll() (int64, error) {
	return service.repo.CountAll()
}

func (service *DocumentService) ProjectDocuments(projectId string) ([]models.Document, error) {
	return service.repo.ProjectDocuments(projectId)
}
