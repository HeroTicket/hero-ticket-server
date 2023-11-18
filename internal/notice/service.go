package notice

import "context"

type Service interface {
	CreateNotice(ctx context.Context, n *Notice) (*Notice, error)
	GetNotice(ctx context.Context, id string) (*Notice, error)
	GetNotices(ctx context.Context, page, limit int) ([]*Notice, error)
	UpdateNotice(ctx context.Context, params *NoticeUpdateParams) error
	DeleteNotice(ctx context.Context, id string) error
}
