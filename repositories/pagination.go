package repositories

type PaginationParams struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type PaginationResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationParams(page, pageSize int) PaginationParams {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPaginationResult(page, pageSize int, total int64) PaginationResult {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}