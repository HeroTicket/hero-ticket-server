package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/service/notice"
	"github.com/heroticket/internal/service/user"
	"go.uber.org/zap"
)

type NoticeCtrl struct {
	notice notice.Service
	user   user.Service
}

func NewNoticeCtrl(notice notice.Service, user user.Service) *NoticeCtrl {
	return &NoticeCtrl{
		notice: notice,
		user:   user,
	}
}

func (c *NoticeCtrl) Pattern() string {
	return "/notices"
}

func (c *NoticeCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.Notices)
	r.Get("/{id}", c.Notice)

	return r
}

// Notices godoc
//
// @Summary Get notices
// @Description Get notices
// @Tags notices
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonResponse
// @Failure 500 {object} CommonResponse
// @Router /notices [get]
func (c *NoticeCtrl) Notices(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var page, limit int64

	if pageStr == "" {
		page = 1
	} else {
		pageInt, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			zap.L().Error("invalid page", zap.Error(err))
			ErrorJSON(w, "invalid page")
			return
		}

		page = pageInt
	}

	if limitStr == "" {
		limit = 10
	} else {
		limitInt, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			zap.L().Error("invalid limit", zap.Error(err))
			ErrorJSON(w, "invalid limit")
			return
		}

		limit = limitInt
	}

	notices, err := c.notice.GetNotices(r.Context(), page, limit)
	if err != nil {
		zap.L().Error("failed to get notices", zap.Error(err))
		ErrorJSON(w, "failed to get notices", http.StatusInternalServerError)
		return
	}

	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "notices retrieved",
		Data:    notices,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}

// Notice godoc
//
// @Summary Get notice
// @Description Get notice
// @Tags notices
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonResponse
// @Failure 500 {object} CommonResponse
// @Router /notices/{id} [get]
func (c *NoticeCtrl) Notice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	n, err := c.notice.GetNotice(r.Context(), id)
	if err != nil {
		zap.L().Error("failed to get notice", zap.Error(err))
		if err == notice.ErrNotFound {
			ErrorJSON(w, "notice not found", http.StatusNotFound)
		} else {
			ErrorJSON(w, "failed to get notice", http.StatusInternalServerError)
		}
		return
	}

	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "notice retrieved",
		Data:    n,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}
