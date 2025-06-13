package update

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
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
		userId, err := validapi.Int64UrlParam(r, "user_id")
		if err != nil {
			slog.Error(err.Error())
			render.JSON(w, r, Response{Response: response.Error(err.Error())})
			return
		}

		noteId, err := validapi.Int64UrlParam(r, "note_id")
		if err != nil {
			slog.Error(err.Error())
			render.JSON(w, r, Response{Response: response.Error(err.Error())})
			return
		}

		var req Request

		err = validapi.JsonBodyDecode(r, &req)
		if err != nil {
			slog.Error(err.Error())
			render.JSON(w, r, Response{Response: response.Error(err.Error())})
			return
		}

		slog.Info(
			"update note request",
			slog.Int64("user_id", userId),
			slog.Int64("note_id", noteId),
			slog.Any("request", req),
		)

		if err := validapi.Request(req); err != nil {
			slog.Error(err.Error())
			render.JSON(w, r, Response{Response: response.Error(err.Error())})
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
