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
