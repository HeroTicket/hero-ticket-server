package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/service/user"
	"go.uber.org/zap"
)

type ProfileCtrl struct {
	user user.Service
}

func NewProfileCtrl() *ProfileCtrl {
	return &ProfileCtrl{}
}

func (c *ProfileCtrl) Pattern() string {
	return "/profile"
}

func (c *ProfileCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/profile/{name}", c.profile)

	return r
}

// Profile godoc
//
//	@Tags			profile
//	@Summary		returns user profile
//	@Description	returns user profile
//	@Accept			json
//	@Produce		json
//	@Param 			name	path	string	true	"user name"
//	@Success		200			{object}	CommonResponse
//	@Failure		400			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Router			/profile	[get]
func (c *ProfileCtrl) profile(w http.ResponseWriter, r *http.Request) {
	// 1. check params
	name := chi.URLParam(r, "name")

	// 2. get user
	_, err := c.user.FindUserByName(r.Context(), name)
	if err != nil {
		zap.L().Error("failed to find user", zap.Error(err))
		if err == user.ErrUserNotFound {
			ErrorJSON(w, "user not found", http.StatusBadRequest)
		} else {
			ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		}
		return
	}

	// 3. get purchased tickets and issued tickets

	// 4. return user profile
}
