package create

import (
	"context"
	"log/slog"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/lib/logger/sl"
)

type Request struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content,omitempty"`
}

type Response struct {
	Id int64 `json:"id,omitempty"`
	response.Response
}

type NotesCreator interface {
	CreateNote(ctx context.Context, userId int64, title string, content string) (int64, error)
}

func New(creator NotesCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		userId, err := validapi.Int64UrlParam(r, "user_id")
		if err != nil {
			slog.Error(err.Error())
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		err = validapi.JsonBodyDecode(r, &req)
		if err != nil {
			slog.Error(err.Error())
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		slog.Info("request body decoded", slog.Any("request", req))

		if err := validapi.Request(req); err != nil {
			slog.Error(err.Error())
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		noteId, err := creator.CreateNote(r.Context(), userId, req.Title, req.Content)
		if err != nil {
			errMsg := "create note error"
			slog.Error(errMsg, sl.Err(err))
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		response.Json(w, r, http.StatusCreated, Response{
			Id:       noteId,
			Response: response.Created(),
		})
	}
}
