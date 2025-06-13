package get

import (
	"context"
	"log/slog"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
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
		userId, err := validapi.Int64UrlParam(r, "user_id")
		if err != nil {
			slog.Error(err.Error())
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		noteId, err := validapi.Int64UrlParam(r, "note_id")
		if err != nil {
			slog.Error(err.Error())
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		slog.Info("get note request", slog.Int64("user_id", userId), slog.Int64("note_id", noteId))

		note, err := noteGetter.GetUserNote(r.Context(), userId, noteId)
		if err != nil {
			errMsg := "get note error"
			slog.Error(errMsg, sl.Err(err))
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		response.Json(w, r, http.StatusOK, Response{
			Note:     &note,
			Response: response.Ok(),
		})
	}
}
