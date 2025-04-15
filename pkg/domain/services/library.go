package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
)

type LibraryService struct {
	repo *repository.LibraryRepository
}

func NewLibraryService(repo *repository.LibraryRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (service *LibraryService) FindStandards(params map[string]any, limit, offset int) ([]models.Project, error) {
	return service.repo.FindStandards(params, limit, offset)
}
