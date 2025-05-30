package data

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalPaginatedRecords, totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	if totalPaginatedRecords == 0 {
		// We return an empty struct because anyway it will be ignored in the JSON
		return Metadata{
			CurrentPage:  page,
			PageSize:     pageSize,
			FirstPage:    1,
			LastPage:     (totalRecords + pageSize - 1) / pageSize,
			TotalRecords: totalRecords,
		}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalPaginatedRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
