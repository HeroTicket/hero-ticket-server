package notice

import "context"

type Service interface {
	CreateNotice(ctx context.Context, n *Notice) (*Notice, error)
	GetNotice(ctx context.Context, id string) (*Notice, error)
	GetNotices(ctx context.Context, page, limit int64) (*Notices, error)
	UpdateNotice(ctx context.Context, params *NoticeUpdateParams) error
	DeleteNotice(ctx context.Context, id int64) error
}

type noticeService struct {
	repo Repository
}

func New(repo Repository) Service {
	return &noticeService{repo: repo}
}

func (svc *noticeService) CreateNotice(ctx context.Context, n *Notice) (*Notice, error) {
	return svc.repo.CreateNotice(ctx, n)
}

func (svc *noticeService) DeleteNotice(ctx context.Context, id int64) error {
	return svc.repo.DeleteNotice(ctx, id)
}

func (svc *noticeService) GetNotice(ctx context.Context, id string) (*Notice, error) {
	return svc.repo.GetNotice(ctx, id)
}

func (svc *noticeService) GetNotices(ctx context.Context, page int64, limit int64) (*Notices, error) {
	return svc.repo.GetNotices(ctx, page, limit)
}

func (svc *noticeService) UpdateNotice(ctx context.Context, params *NoticeUpdateParams) error {
	return svc.repo.UpdateNotice(ctx, params)
}
