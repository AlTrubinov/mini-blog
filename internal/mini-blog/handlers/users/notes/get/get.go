package get

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/logger/sl"
	"mini-blog/internal/models/note"
)

type Response struct {
	Note *note.Note `json:"note,omitempty"`
	response.Response
}

type NoteGetter interface {
	GetUserNote(ctx context.Context, userId int64, noteId int64) (note.Note, error)
}

func New(noteGetter NoteGetter) http.HandlerFunc {
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

		slog.Info("get note request", slog.Int64("user_id", userId), slog.Int64("note_id", noteId))

		note, err := noteGetter.GetUserNote(r.Context(), userId, noteId)
		if err != nil {
			errMsg := "get note error"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		render.JSON(w, r, Response{
			Note:     &note,
			Response: response.Ok(),
		})
	}
}
