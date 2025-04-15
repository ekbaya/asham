package services

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type LibraryService struct {
	repo *repository.LibraryRepository
}

func NewLibraryService(repo *repository.LibraryRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (service *LibraryService) FindStandards(params map[string]any, limit, offset int) ([]map[string]any, error) {
	var standards []map[string]any

	projects, err := service.repo.FindStandards(params, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(projects) > 0 {
		for _, project := range projects {
			pageCount := 20 
			if project.Standard != nil && project.Standard.FileURL != "" {
				calculatedPages, err := service.getPDFPageCount(project.Standard.FileURL)
				if err == nil {
					pageCount = calculatedPages
				} else {
					log.Printf("Error calculating PDF pages for standard ID %v: %v", project.ID, err)
				}
			}

			standard := map[string]any{
				"id":             project.ID,
				"title":          project.Title,
				"description":    project.Description,
				"sector":         project.Sector,
				"committee":      project.TechnicalCommittee.Code,
				"language":       "English",
				"published_date": project.PublishedDate,
				"pages":          pageCount,
			}
			standards = append(standards, standard)
		}
		return standards, nil
	} else {
		return nil, err
	}
}

// getPDFPageCount downloads the PDF from the URL and calculates the number of pages
func (service *LibraryService) getPDFPageCount(pdfURL string) (int, error) {
	// Parse the URL
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return 0, err
	}

	// Handle local file path (for assets directory)
	if parsedURL.Scheme == "" || parsedURL.Scheme == "file" {
		// Assuming this is a local file path
		return api.PageCountFile(pdfURL)
	}

	// For remote URLs, download the file temporarily
	resp, err := http.Get(pdfURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Create a temporary file to store the PDF
	tmpFile, err := os.CreateTemp("", "standard-*.pdf")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Copy the PDF data to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		return 0, err
	}
	tmpFile.Close()

	// Get the page count from the downloaded file
	pageCount, err := api.PageCountFile(tmpFile.Name())
	if err != nil {
		return 0, err
	}

	return pageCount, nil
}
