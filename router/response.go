package router

import (
	"math"
	"runtime"

	"github.com/enesanbar/go-service/errors"
)

type PagedResponse struct {
	TotalItems  int64       `json:"total_items"`
	NumOfPages  int64       `json:"num_of_pages"`
	CurrentPage int64       `json:"current_page"`
	PageSize    int64       `json:"page_size"`
	PrevPage    *int64      `json:"prev_page,omitempty"`
	NextPage    *int64      `json:"next_page,omitempty"`
	Items       interface{} `json:"items"`
}

func NewPagedResponse(items interface{}, page int64, pageSize int64, count int64) *PagedResponse {
	p := &PagedResponse{}
	numOfPages := int64(math.Ceil(float64(count) / float64(pageSize)))

	if page != numOfPages && page <= numOfPages {
		nextPage := page + 1
		p.NextPage = &nextPage
	}

	if page != 1 && page == numOfPages {
		prevPage := page - 1
		p.PrevPage = &prevPage
	}

	p.Items = items
	p.TotalItems = count
	p.NumOfPages = numOfPages
	p.CurrentPage = page
	p.PageSize = pageSize

	return p
}

type ApiResponse struct {
	Status int         `json:"status"`
	Err    error       `json:"-"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Line   int         `json:"-"`
	File   string      `json:"-"`
}

func NewApiResponse(code int, data interface{}, err error) ApiResponse {
	response := ApiResponse{Err: err, Data: data, Status: code}

	if e, ok := err.(errors.Error); ok && e.Err != nil {
		_, f, n, _ := runtime.Caller(1)
		response.Error = errors.ErrorMessage(err)
		response.File = f
		response.Line = n
	}

	return response
}
