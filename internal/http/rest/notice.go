package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/jwt"
	"github.com/heroticket/internal/notice"
	"github.com/heroticket/internal/user"
)

type noticeCtrl struct {
	jwt    jwt.Service
	notice notice.Service
	user   user.Service
	tx     db.Tx
}

func NewNoticeCtrl(jwt jwt.Service, notice notice.Service, user user.Service, tx db.Tx) *noticeCtrl {
	return &noticeCtrl{
		jwt:    jwt,
		notice: notice,
		user:   user,
		tx:     tx,
	}
}

func (c *noticeCtrl) Pattern() string {
	return "/notices"
}

func (c *noticeCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", nil)
	r.Get("/{id}", nil)

	r.Group(func(r chi.Router) {
		r.Use(AccessTokenRequired(c.jwt))

		r.Post("/", nil)
		r.Put("/{id}", nil)
	})

	return r
}
