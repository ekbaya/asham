package models

type SharepointDocument struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	WebURL       string `json:"webUrl"`
	CreatedBy    string `json:"createdBy"`
	LastModified string `json:"lastModified"`
	EmbedUrl     string `json:"embedUrl"`
}
