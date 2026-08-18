package model

type Page[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func NewPage[T any](items []T, total, offset, limit int) Page[T] {
	if items == nil {
		items = []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	return Page[T]{Items: items, Total: total, Offset: offset, Limit: limit}
}
