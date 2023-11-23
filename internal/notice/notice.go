package notice

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("notice not found")
	ErrNothingToUpdate = errors.New("nothing to update")
)

type Notice struct {
	ID        string    `json:"id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type NoticeUpdateParams struct {
	ID      string
	Title   string
	Content string
}

type Pagination struct {
	Total       int64 `json:"total"`
	Pages       int64 `json:"pages"`
	CurrentPage int64 `json:"current_page"`
	Limit       int64 `json:"limit"`
	Start       int64 `json:"start"`
	End         int64 `json:"end"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

type Notices struct {
	Items      []*Notice   `json:"items"`
	Pagination *Pagination `json:"pagination"`
}
