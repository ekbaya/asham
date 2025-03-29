package utilities

// PaginationData represents pagination information for API responses
type PaginationData struct {
	CurrentPage  int  `json:"current_page"`
	TotalPages   int  `json:"total_pages"`
	TotalItems   int  `json:"total_items"`
	ItemsPerPage int  `json:"items_per_page"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

// GeneratePaginationData creates pagination metadata based on the current request
func GeneratePaginationData(limit int, currentPage int, totalItems int) PaginationData {
	totalPages := 0
	if limit > 0 {
		totalPages = (totalItems + limit - 1) / limit // Ceiling division
	}

	// Make sure we don't have a current page higher than total pages
	if currentPage > totalPages && totalPages > 0 {
		currentPage = totalPages
	}

	return PaginationData{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
		HasNext:      currentPage < totalPages,
		HasPrevious:  currentPage > 1,
	}
}
