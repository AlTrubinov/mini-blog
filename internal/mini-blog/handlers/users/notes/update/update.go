package update

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/logger/sl"
)

type Request struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content,omitempty"`
}

type Response struct {
	response.Response
}

type NoteUpdater interface {
	UpdateNote(ctx context.Context, userId int64, noteId int64, title string, content string) error
}

func New(noteUpdater NoteUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
		if err != nil {
			errMsg := "invalid user id"
			slog.Info(errMsg)

			render.JSON(w, r, Response{
				Response: response.Error(errMsg),
			})
			return
		}

		noteId, err := strconv.ParseInt(chi.URLParam(r, "note_id"), 10, 64)
		if err != nil {
			errMsg := "invalid note id"
			slog.Info(errMsg)

			render.JSON(w, r, Response{
				Response: response.Error(errMsg),
			})
			return
		}

		var req Request

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			errMsg := "decode request body failed"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		slog.Info(
			"update note request",
			slog.Int64("user_id", userId),
			slog.Int64("note_id", noteId),
			slog.Any("request", req),
		)

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			errMsg := "request validation failed"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.ValidationError(validateErr)})
			return
		}

		err = noteUpdater.UpdateNote(r.Context(), userId, noteId, req.Title, req.Content)
		if err != nil {
			errMsg := "update note error"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		render.JSON(w, r, Response{Response: response.Ok()})
	}
}
