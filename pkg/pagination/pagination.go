package pagination

import (
	"math"
	"net/http"
	"strconv"
)

// Pagination holds pagination details.
type Pagination struct {
	TotalPages  int   `json:"totalPages"`
	TotalCount  int64 `json:"totalCount"`
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
}

// New creates a new Pagination instance.
func New(currentPage, pageSize, totalCount int) *Pagination {
	if currentPage <= 0 {
		currentPage = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &Pagination{
		TotalPages:  totalPages,
		TotalCount:  int64(totalCount),
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}
}

// GetPaginationParams extracts page and pageSize from the request.
func GetPaginationParams(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	return page, pageSize
}
