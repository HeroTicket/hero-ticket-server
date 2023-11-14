package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TicketCtrl struct {
}

func NewTicketCtrl() *TicketCtrl {
	return &TicketCtrl{}
}

func (c *TicketCtrl) Pattern() string {
	return "/tickets"
}

func (c *TicketCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Post("/create", nil)
	r.Get("/create-qr", nil)
	r.Post("/create-callback", nil)
	r.Get("/list", nil)
	r.Get("/list/{id}", nil)
	r.Get("/{id}/purchase-qr", nil)
	r.Post("/{id}/purchase-callback", nil)
	r.Get("/{id}/verify-qr", nil)
	r.Post("/{id}/verify-callback", nil)

	return r
}
