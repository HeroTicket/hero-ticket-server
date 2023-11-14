package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserCtrl struct {
}

func NewUserCtrl() *UserCtrl {
	return &UserCtrl{}
}

func (ctrl *UserCtrl) Pattern() string {
	return "/users"
}

func (ctrl *UserCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/login-qr", nil)
	r.Post("/login-callback", nil)
	r.Post("/logout", nil)
	r.Post("/refresh-token", nil)
	r.Post("/create-tba", nil)
	r.Get("/purchased-tickets", nil)
	r.Get("/issued-tickets", nil)

	return r
}
