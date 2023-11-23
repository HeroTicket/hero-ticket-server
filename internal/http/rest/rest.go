package rest

import "github.com/heroticket/internal/notice"

type CommonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type NoticesResponse struct {
	Notices    []*notice.Notice   `json:"notices"`
	Pagination *notice.Pagination `json:"pagination"`
}
