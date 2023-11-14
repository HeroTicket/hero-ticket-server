package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/heroticket/internal/ws"
)

type Controller interface {
	Pattern() string
	Handler() http.Handler
}

type router struct {
	*chi.Mux
}

func newRouter(version string, ctrls ...Controller) *router {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route(fmt.Sprintf("/%s", version), func(r chi.Router) {
			for _, ctrl := range ctrls {
				r.Mount(ctrl.Pattern(), ctrl.Handler())
			}
		})
	})

	r.HandleFunc("/ws", ws.Serve())

	return &router{r}
}