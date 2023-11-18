package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/jwt"
	"github.com/heroticket/internal/notice"
	"github.com/heroticket/internal/user"
	"go.uber.org/zap"
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

	r.Get("/", c.Notices)
	r.Get("/{id}", c.Notice)

	r.Group(func(r chi.Router) {
		r.Use(AccessTokenRequired(c.jwt))

		r.Post("/", c.CreateNotice)
		r.Put("/{id}", c.UpdateNotice)
	})

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
// @Success 200 {object} NoticesResponse
// @Failure 400 {object} CommonResponse
// @Failure 500 {object} CommonResponse
// @Router /notices [get]
func (c *noticeCtrl) Notices(w http.ResponseWriter, r *http.Request) {
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

	notices, pagination, err := c.notice.GetNotices(r.Context(), page, limit)
	if err != nil {
		zap.L().Error("failed to get notices", zap.Error(err))
		ErrorJSON(w, "failed to get notices", http.StatusInternalServerError)
		return
	}

	var resp NoticesResponse

	resp.Notices = notices
	resp.Pagination = pagination

	_ = WriteJSON(w, http.StatusOK, resp)
}

// Notice godoc
//
// @Summary Get notice
// @Description Get notice
// @Tags notices
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} NoticeResponse
// @Failure 400 {object} CommonResponse
// @Failure 500 {object} CommonResponse
// @Router /notices/{id} [get]
func (c *noticeCtrl) Notice(w http.ResponseWriter, r *http.Request) {
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

	_ = WriteJSON(w, http.StatusOK, n)
}

// CreateNotice godoc
//
// @Summary Create notice
// @Description Create notice
// @Tags notices
// @Produce json
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Success 201 {object} CommonResponse
// @Failure 400 {object} CommonResponse
// @Failure 401 {object} CommonResponse
// @Failure 500 {object} CommonResponse
// @Router /notices [post]
func (c *noticeCtrl) CreateNotice(w http.ResponseWriter, r *http.Request) {
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("user not found", zap.Error(err))
		ErrorJSON(w, "user not found", http.StatusUnauthorized)
		return
	}

	u, err := c.user.FindUserByDID(r.Context(), jwtUser.DID)
	if err != nil {
		zap.L().Error("failed to find user", zap.Error(err))
		ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		return
	}

	if !u.IsAdmin {
		ErrorJSON(w, "user is not admin", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" {
		ErrorJSON(w, "title is required", http.StatusBadRequest)
		return
	}

	if content == "" {
		ErrorJSON(w, "content is required", http.StatusBadRequest)
		return
	}

	n := &notice.Notice{
		Title:   title,
		Content: content,
	}

	_, err = c.notice.CreateNotice(r.Context(), n)
	if err != nil {
		zap.L().Error("failed to create notice", zap.Error(err))
		ErrorJSON(w, "failed to create notice", http.StatusInternalServerError)
		return
	}

	var resp CommonResponse

	resp.Status = http.StatusCreated
	resp.Message = "notice created"

	_ = WriteJSON(w, http.StatusCreated, resp)
}

// UpdateNotice godoc
//
// @Summary Update notice
// @Description Update notice
// @Tags notices
// @Produce json
// @Param id path string true "id"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonResponse
// @Failure 401 {object} CommonResponse
// @Failure 500 {object} CommonResponse
// @Router /notices/{id} [put]
func (c *noticeCtrl) UpdateNotice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("user not found", zap.Error(err))
		ErrorJSON(w, "user not found", http.StatusUnauthorized)
		return
	}

	u, err := c.user.FindUserByDID(r.Context(), jwtUser.DID)
	if err != nil {
		zap.L().Error("failed to find user", zap.Error(err))
		ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		return
	}

	if !u.IsAdmin {
		ErrorJSON(w, "user is not admin", http.StatusForbidden)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" {
		ErrorJSON(w, "title is required", http.StatusBadRequest)
		return
	}

	if content == "" {
		ErrorJSON(w, "content is required", http.StatusBadRequest)
		return
	}

	params := &notice.NoticeUpdateParams{
		ID:      id,
		Title:   title,
		Content: content,
	}

	err = c.notice.UpdateNotice(r.Context(), params)
	if err != nil {
		zap.L().Error("failed to update notice", zap.Error(err))
		ErrorJSON(w, "failed to update notice", http.StatusInternalServerError)
		return
	}

	var resp CommonResponse

	resp.Status = http.StatusOK
	resp.Message = "notice updated"

	_ = WriteJSON(w, http.StatusOK, resp)
}
