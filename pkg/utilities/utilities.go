package utilities

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func GetPDFPageCount(pdfURL string) (int, error) {
	// Parse the URL
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return 0, err
	}

	// Handle local file path (for assets directory)
	if parsedURL.Scheme == "" || parsedURL.Scheme == "file" {
		// Assuming this is a local file path
		basePath := "/home/ubuntu/projects/asham"
		if pdfURL[0] == '/' { // Check if it starts with a forward slash
			pdfURL = basePath + pdfURL // Preserve the leading "/"
		} else {
			pdfURL = basePath + "/" + pdfURL
		}
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

func ToUpperUnderscore(s string) string {
	// Replace spaces and hyphens with underscores
	re := regexp.MustCompile(`[ -]+`)
	s = re.ReplaceAllString(s, "_")
	return strings.ToUpper(s)
}
