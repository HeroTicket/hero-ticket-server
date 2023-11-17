package notice

import "time"

type Notice struct {
	ID        string    `json:"id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type MediumType string

const (
	MediumTypeImage MediumType = "image"
	MediumTypeVideo MediumType = "video"
)

type Medium struct {
	ID         string     `json:"id" bson:"_id"`
	NoticeID   string     `json:"notice_id" bson:"notice_id"`
	MediumType MediumType `json:"medium_type" bson:"medium_type"`
	Content    string     `json:"content" bson:"content"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
}
