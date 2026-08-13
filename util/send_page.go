package util

import "net/http"

type PaginatedData struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int64 `json:"totalPages"`
}

func SendPage(w http.ResponseWriter, data any, page, limit, cnt int64) {
	var totalPages int64
	if limit > 0 {
		totalPages = (cnt + limit - 1) / limit
	}

	paginatedData := PaginatedData{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: cnt,
			TotalPages: totalPages,
		},
	}

	SendData(w, http.StatusOK, paginatedData)
}
