package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/service/user"
)

type ProfileCtrl struct {
	ticket ticket.Service
	user   user.Service
}

func NewProfileCtrl(ticket ticket.Service, user user.Service) *ProfileCtrl {
	return &ProfileCtrl{
		ticket: ticket,
		user:   user,
	}
}

func (c *ProfileCtrl) Pattern() string {
	return "/profile"
}

func (c *ProfileCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/{name}", c.profile)

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
//	@Router			/v1/profile/{name} [get]
func (c *ProfileCtrl) profile(w http.ResponseWriter, r *http.Request) {
	// 1. check params
	name := chi.URLParam(r, "name")

	// 2. get user
	u, err := c.user.FindUserByName(r.Context(), name)
	if err != nil {
		logger.Error("failed to find user", "error", err)
		if err == user.ErrUserNotFound {
			ErrorJSON(w, "user not found", http.StatusBadRequest)
		} else {
			ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		}
		return
	}

	// TODO: get user profile
	// 3. get purchased tickets and issued tickets from db

	// 4. return user profile
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully get user profile",
		Data:    u,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}
