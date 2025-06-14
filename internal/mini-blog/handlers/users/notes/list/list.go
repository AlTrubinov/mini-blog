package list

import (
	"context"
	"log/slog"
	"mini-blog/internal/lib/auth"
	"net/http"
	"strconv"
	"strings"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/models/note"
	"mini-blog/pkg/apperror"
)

type Response struct {
	Notes []note.Note `json:"notes"`
	response.Response
}

type NotesList interface {
	GetUserNotes(ctx context.Context, userId int64, limit int, offset int, order string) ([]note.Note, error)
}

const (
	defaultParamLimit  = 10
	defaultParamOffset = 0

	orderByAsc  = "ASC"
	orderByDesc = "DESC"
)

func New(list NotesList, tm *auth.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := validapi.Int64UrlParam(r, "user_id")
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		if err = tm.CheckUserAccess(r.Context(), userId); err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		query := r.URL.Query()
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			limit = defaultParamLimit
		}
		offset, err := strconv.Atoi(query.Get("offset"))
		if err != nil {
			offset = defaultParamOffset
		}
		order := strings.ToUpper(query.Get("order"))
		if order != orderByAsc && order != orderByDesc {
			order = orderByAsc
		}

		slog.Info("query params parsed", slog.Int64("user_id", userId), slog.Int("limit", limit), slog.Int("offset", offset), slog.String("order", order))

		notes, err := list.GetUserNotes(r.Context(), userId, limit, offset, order)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		response.Json(w, r, http.StatusOK, Response{
			Notes:    notes,
			Response: response.Ok(),
		})
	}
}
