package app

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/heroticket/internal/app/ws"
)

//go:embed icon.png
var res embed.FS

type Controller interface {
	Pattern() string
	Handler() http.Handler
}

type router struct {
	*chi.Mux
}

func newRouter(version string, ctrls ...Controller) *router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Route(fmt.Sprintf("/%s", version), func(r chi.Router) {
		for _, ctrl := range ctrls {
			r.Mount(ctrl.Pattern(), ctrl.Handler())
		}
	})

	r.Get("/", nameHandler)
	r.Get("/status", statusHandler)
	r.Get("/favicon.ico", faviconHandler)

	r.HandleFunc("/ws", ws.Serve())

	return &router{r}
}

// Name godoc
//
// @Summary Get name
// @Description returns name
// @Tags common
// @Accept plain
// @Produce plain
// @Success 200 {string} string "Hero Ticket"
// @Router / [get]
func nameHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hero Ticket"))
}

// Status godoc
//
// @Summary Get status
// @Description returns status
// @Tags common
// @Accept plain
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /status [get]
func statusHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}

// Favicon godoc
//
// @Summary Get favicon
// @Description returns favicon
// @Tags common
// @Accept plain
// @Produce plain
// @Success 200 {file} file "favicon.ico"
// @Router /favicon.ico [get]
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	f, err := res.ReadFile("icon.png")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(f)
}
