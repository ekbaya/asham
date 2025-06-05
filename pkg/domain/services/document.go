package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ekbaya/asham/pkg/config"
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type DocumentService struct {
	repo         *repository.DocumentRepository
	projectRepo  *repository.ProjectRepository
	client       *msgraphsdk.GraphServiceClient
	tokenManager *TokenManager
}

func NewDocumentService(repo *repository.DocumentRepository, projectRepo *repository.ProjectRepository, client *msgraphsdk.GraphServiceClient, tokenManager *TokenManager) *DocumentService {
	return &DocumentService{repo: repo, projectRepo: projectRepo, client: client, tokenManager: tokenManager}
}

func (service *DocumentService) Create(doc *models.Document) error {
	doc.ID = uuid.New()
	doc.CreatedAt = time.Now()
	return service.repo.Create(doc)
}

func (service *DocumentService) UploadStandard(doc *models.Document, project *models.Project) error {
	doc.ID = uuid.New()
	doc.CreatedAt = time.Now()
	project.ID = uuid.New()
	doc.CreatedAt = time.Now()
	return service.repo.UploadStandard(doc, project)
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
			if err := service.projectRepo.UpdateProject(&project); err != nil {
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

func (service *DocumentService) UpdateProjectRelatedDoc(projectId, docTitle, docRef, docDescription, fileURL, member string) error {
	return service.repo.UpdateProjectRelatedDoc(projectId, docTitle, docRef, docDescription, fileURL, member)
}

func (service *DocumentService) UpdateMeetingMinutes(meetingId, fileURL, member string) error {
	return service.repo.UpdateMeetingMinutes(meetingId, fileURL, member)
}

func (service *DocumentService) ListDocuments(ctx context.Context) ([]models.SharepointDocument, error) {
	// Retrieve the token from the token manager
	userToken, err := service.tokenManager.RetrieveToken(ctx)
	if err != nil {
		return nil, err
	}

	userEmail := config.GetConfig().AZURE_USER_EMAIL

	// Safely encode email and search query
	safeEmail := url.PathEscape(userEmail)
	searchQuery := "*.docx"
	escapedQuery := url.QueryEscape(searchQuery) // becomes %2A.docx

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/drive/root/search(q=%s)", safeEmail, escapedQuery)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch items: %s, body: %s", resp.Status, string(body))
	}

	var result struct {
		Value []struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			WebUrl string `json:"webUrl"`
			File   *struct {
				MimeType string `json:"mimeType"`
			} `json:"file"`
			CreatedBy struct {
				User struct {
					DisplayName string `json:"displayName"`
				} `json:"user"`
			} `json:"createdBy"`
			LastModifiedDateTime string `json:"lastModifiedDateTime"`
		} `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var documents []models.SharepointDocument
	for _, item := range result.Value {
		if item.File != nil && item.File.MimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
			documents = append(documents, models.SharepointDocument{
				ID:           item.Id,
				Name:         item.Name,
				WebURL:       item.WebUrl,
				CreatedBy:    item.CreatedBy.User.DisplayName,
				LastModified: item.LastModifiedDateTime,
			})
		}
	}

	return documents, nil
}

func (service *DocumentService) GetDocument(ctx context.Context, documentId string) (*models.SharepointDocument, error) {
	// Retrieve the token from the token manager
	token, err := service.tokenManager.RetrieveToken(ctx)
	if err != nil {
		return nil, err
	}

	userEmail := config.GetConfig().AZURE_USER_EMAIL

	oneDriveUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/drive/items/%s", userEmail, documentId)
	req, _ := http.NewRequestWithContext(ctx, "GET", oneDriveUrl, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document from OneDrive: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch document details: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		WebUrl string `json:"webUrl"`
		File   *struct {
			MimeType string `json:"mimeType"`
		} `json:"file"`
		CreatedBy struct {
			User struct {
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"createdBy"`
		LastModifiedDateTime string `json:"lastModifiedDateTime"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode OneDrive document response: %v", err)
	}
	if result.File == nil {
		return nil, fmt.Errorf("item is not a file")
	}
	if result.File.MimeType != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		return nil, fmt.Errorf("item is not a Word document (mimeType: %s)", result.File.MimeType)
	}

	// Call the /preview endpoint to get the embedUrl
	previewUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/preview", documentId)
	previewReq, _ := http.NewRequestWithContext(ctx, "POST", previewUrl, nil)
	previewReq.Header.Set("Authorization", "Bearer "+token)
	previewReq.Header.Set("Accept", "application/json")

	embedUrl := ""
	previewResp, err := http.DefaultClient.Do(previewReq)
	if err == nil && previewResp.StatusCode == 200 {
		var previewResult struct {
			GetUrl string `json:"getUrl"`
		}
		if err := json.NewDecoder(previewResp.Body).Decode(&previewResult); err == nil {
			embedUrl = previewResult.GetUrl
		}
	}
	if previewResp != nil {
		previewResp.Body.Close()
	}

	return &models.SharepointDocument{
		ID:           result.Id,
		Name:         result.Name,
		WebURL:       result.WebUrl,
		CreatedBy:    result.CreatedBy.User.DisplayName,
		LastModified: result.LastModifiedDateTime,
		EmbedUrl:     embedUrl,
	}, nil
}
