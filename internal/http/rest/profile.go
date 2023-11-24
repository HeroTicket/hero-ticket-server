package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ProfileCtrl struct {
}

func NewProfileCtrl() *ProfileCtrl {
	return &ProfileCtrl{}
}

func (c *ProfileCtrl) Pattern() string {
	return "/profile"
}

func (c *ProfileCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/profile", c.profile)

	return r
}

// Profile godoc
//
//	@Tags			profile
//	@Summary		returns user profile
//	@Description	returns user profile
//	@Produce		json
//	@Param			id			query		string	false	"id"
//	@Param			name		query		string	false	"name"
//	@Success		200			{object}	CommonResponse
//	@Failure		400			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Router			/users/profile [get]
func (c *ProfileCtrl) profile(w http.ResponseWriter, r *http.Request) {
	// 1. check params
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")

	if id == "" && name == "" {
		ErrorJSON(w, "id or name is required", http.StatusBadRequest)
		return
	}

	// 2. get user profile

	// 3. return user profile
}
