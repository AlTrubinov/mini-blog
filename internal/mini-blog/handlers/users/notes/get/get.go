package get

import (
	"context"
	"log/slog"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/models/note"
	"mini-blog/pkg/apperror"
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
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		noteId, err := validapi.Int64UrlParam(r, "note_id")
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		slog.Info("get note request", slog.Int64("user_id", userId), slog.Int64("note_id", noteId))

		note, err := noteGetter.GetUserNote(r.Context(), userId, noteId)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		response.Json(w, r, http.StatusOK, Response{
			Note:     &note,
			Response: response.Ok(),
		})
	}
}
