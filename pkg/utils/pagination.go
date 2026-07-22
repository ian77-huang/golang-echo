package utils

type PaginationPage struct {
	Value         int  `json:"value"`
	IsCurrentPage bool `json:"is_current_page"`
}

type Pagination struct {
	CurrentPage int              `json:"current_page"`
	TotalPages  int              `json:"total_pages"`
	Pages       []PaginationPage `json:"pages"`
}

func GetPagination(currentPage, totalPages, spacing int) *Pagination {
	if totalPages <= 0 {
		return &Pagination{
			CurrentPage: currentPage,
			TotalPages:  0,
			Pages:       []PaginationPage{},
		}
	}

	if spacing <= 0 {
		spacing = 5
	}

	half := spacing / 2
	start := currentPage - half
	end := start + spacing - 1

	if start < 1 {
		start = 1
		end = spacing
	}

	if end > totalPages {
		end = totalPages
		start = end - spacing + 1
		if start < 1 {
			start = 1
		}
	}

	pages := make([]PaginationPage, 0, end-start+1)
	for i := start; i <= end; i++ {
		pages = append(pages, PaginationPage{
			Value:         i,
			IsCurrentPage: i == currentPage,
		})
	}

	return &Pagination{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Pages:       pages,
	}
}
