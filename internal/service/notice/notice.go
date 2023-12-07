package notice

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("notice not found")
	ErrNothingToUpdate = errors.New("nothing to update")
)

type Notice struct {
	ID        int64  `json:"id" bson:"_id"`
	Title     string `json:"title" bson:"title"`
	Content   string `json:"content" bson:"content"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64  `json:"updatedAt" bson:"updatedAt"`
}

type NoticeUpdateParams struct {
	ID      int64
	Title   string
	Content string
}

type Pagination struct {
	Total       int64 `json:"total"`
	Pages       int64 `json:"pages"`
	CurrentPage int64 `json:"currentPage"`
	Limit       int64 `json:"limit"`
	Start       int64 `json:"start"`
	End         int64 `json:"end"`
	HasNext     bool  `json:"hasNext"`
	HasPrev     bool  `json:"hasPrev"`
}

type Notices struct {
	Items      []*Notice   `json:"items"`
	Pagination *Pagination `json:"pagination"`
}
