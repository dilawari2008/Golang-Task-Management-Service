package models

// Pagination represents pagination metadata
type Pagination struct {
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// NewPagination creates a new pagination struct
func NewPagination(totalCount int64, page, limit int) *Pagination {
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return &Pagination{
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}