package helpers

import (
	"math"
	"reflect"
)

type Pagination interface {
	GetPage() int
	GetLimit() int
	GetTotal() int
	GetPages() int
	GetData() interface{}
}

type pagination struct {
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
	Pages   int         `json:"pages"`
	Data    interface{} `json:"data"`
}

func (p *pagination) GetPage() int {
	return p.Page
}

func (p *pagination) GetLimit() int {
	return p.Limit
}

func (p *pagination) GetTotal() int {
	return p.Total
}

func (p *pagination) GetPages() int {
	return p.Pages
}

func (p *pagination) GetData() interface{} {
	return p.Data
}

func Paginate(list interface{}, page int, limit int) Pagination {
	total := reflect.ValueOf(list).Len()
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}

	if total == 0 {
		return &pagination{
			Page:    1,
			Limit:   0,
			Total:   0,
			Pages:   1,
			Data:    []interface{}{},
		}
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))
	if page > pages {
		page = pages
	}
	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}
	data := reflect.ValueOf(list).Slice(start, end).Interface()

	return &pagination{
		Page:    page,
		Limit:   limit,
		Total:   total,
		Pages:   pages,
		Data:    data,
	}
}
