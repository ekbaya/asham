package collaboration

import (
	"context"

	"github.com/ekbaya/asham/pkg/config"
	"golang.org/x/oauth2"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func CreateGoogleDoc(title string) (string, error) {
	config := config.GetConfig()
	ctx := context.Background()
	srv, err := docs.NewService(ctx, option.WithTokenSource(
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GOOGLE_CLIENT_TOKEN}),
	))
	if err != nil {
		return "", err
	}

	doc, err := srv.Documents.Create(&docs.Document{
		Title: title,
	}).Do()

	if err != nil {
		return "", err
	}

	return doc.DocumentId, nil
}
