package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (service *CommentService) Create(comment *models.CommentObservation) error {
	return service.repo.Create(comment)
}

func (service *CommentService) GetByID(id uuid.UUID) (*models.CommentObservation, error) {
	return service.repo.GetByID(id)
}

func (service *CommentService) GetAll() ([]models.CommentObservation, error) {
	return service.repo.GetAll()
}

func (service *CommentService) Update(comment *models.CommentObservation) error {
	return service.repo.Update(comment)
}

func (service *CommentService) Delete(id uuid.UUID) error {
	return service.repo.Delete(id)
}

func (service *CommentService) GetByProjectID(projectID uuid.UUID) ([]models.CommentObservation, error) {
	return service.repo.GetByProjectID(projectID)
}
