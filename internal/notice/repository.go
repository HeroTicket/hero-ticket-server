package notice

import "context"

type Command interface {
	CreateNotice(ctx context.Context, n *Notice) (*Notice, error)
	UpdateNotice(ctx context.Context, id string, n *Notice) (*Notice, error)
	DeleteNotice(ctx context.Context, id string) error

	CreateMedium(ctx context.Context, m *Medium) (*Medium, error)
	UpdateMedium(ctx context.Context, id string, m *Medium) (*Medium, error)
	DeleteMedium(ctx context.Context, id string) error
	DeleteMediumsByNoticeID(ctx context.Context, noticeID string) error
}

type Query interface {
	GetNotice(ctx context.Context, id string) (*Notice, error)
	GetNotices(ctx context.Context, page, limit int) ([]*Notice, error)

	GetMedium(ctx context.Context, id string) (*Medium, error)
	GetMediumsByNoticeID(ctx context.Context, noticeID string) ([]*Medium, error)
}

type Repository interface {
	Command
	Query
}
