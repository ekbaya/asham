package utilities

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func GetPDFPageCount(pdfURL string) (int, error) {
	// Parse the URL
	parsedURL, err := url.Parse(strings.TrimPrefix(pdfURL, "/"))
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
