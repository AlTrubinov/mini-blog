package create

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
	UserId  int64  `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
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

		userId, err := strconv.Atoi(
			chi.URLParam(r, "user_id"),
		)
		if err != nil {
			errMsg := "invalid user id"
			slog.Info(errMsg)

			render.JSON(w, r, Response{
				Response: response.Error(errMsg),
			})
			return
		}
		req.UserId = int64(userId)

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			errMsg := "decode request body failed"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		slog.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			errMsg := "request validation failed"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.ValidationError(validateErr)})
			return
		}

		noteId, err := creator.CreateNote(r.Context(), req.UserId, req.Title, req.Content)
		if err != nil {
			errMsg := "create note error"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		render.JSON(w, r, Response{
			Id:       noteId,
			Response: response.Ok(),
		})
	}
}
