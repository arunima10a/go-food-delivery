package utils

import (
	 "math"
)

type Pagination struct {
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int`json:"totalPages"`
	Items interface{} `json:"items"`
}

func NewPagination(page, pageSize int, totalItems int64, items interface{}) *Pagination {
	return &Pagination {
		Page: page,
		PageSize: pageSize,
		TotalItems: totalItems,
		TotalPages: int(math.Ceil(float64(totalItems)/ float64(pageSize))),
		Items: items,

	}
}