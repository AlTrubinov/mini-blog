package list

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/lib/logger/sl"
	"mini-blog/internal/models/note"
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

func New(list NotesList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := validapi.Int64UrlParam(r, "user_id")
		if err != nil {
			slog.Error(err.Error())
			render.JSON(w, r, Response{Response: response.Error(err.Error())})
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
			errMsg := "get notes list error"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		render.JSON(w, r, Response{
			Notes:    notes,
			Response: response.Ok(),
		})
	}
}
