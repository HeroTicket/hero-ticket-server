package notice

import "context"

type Service interface {
	CreateNotice(ctx context.Context, n *Notice) (*Notice, error)
	GetNotice(ctx context.Context, id string) (*Notice, error)
	GetNotices(ctx context.Context, page, limit int) ([]*Notice, error)
	UpdateNotice(ctx context.Context, id string, n *Notice) (*Notice, error)
	DeleteNotice(ctx context.Context, id string) error

	CreateMedium(ctx context.Context, m *Medium) (*Medium, error)
	GetMedium(ctx context.Context, id string) (*Medium, error)
	GetMediumsByNoticeID(ctx context.Context, noticeID string) ([]*Medium, error)
	UpdateMedium(ctx context.Context, id string, m *Medium) (*Medium, error)
	DeleteMedium(ctx context.Context, id string) error
	DeleteMediumsByNoticeID(ctx context.Context, noticeID string) error
}
