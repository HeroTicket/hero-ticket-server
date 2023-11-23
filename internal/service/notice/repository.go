package notice

import "context"

type Command interface {
	CreateNotice(ctx context.Context, n *Notice) (*Notice, error)
	UpdateNotice(ctx context.Context, params *NoticeUpdateParams) error
	DeleteNotice(ctx context.Context, id string) error
}

type Query interface {
	GetNotice(ctx context.Context, id string) (*Notice, error)
	GetNotices(ctx context.Context, page, limit int64) (*Notices, error)
}

type Repository interface {
	Command
	Query
}
