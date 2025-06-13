package delete

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/lib/logger/sl"
)

type Response struct {
	response.Response
}

type NoteDeleter interface {
	DeleteNote(ctx context.Context, userId int64, noteId int64) error
}

func New(noteDeleter NoteDeleter) http.HandlerFunc {
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

		slog.Info("delete note request", slog.Int64("user_id", userId), slog.Int64("note_id", noteId))

		err = noteDeleter.DeleteNote(r.Context(), userId, noteId)
		if err != nil {
			errMsg := "delete note error"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		render.JSON(w, r, Response{Response: response.Ok()})
	}
}
