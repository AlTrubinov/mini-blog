package registration

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/logger/sl"
)

type Request struct {
	Username string `json:"username" validate:"required"`
}

type Response struct {
	Id       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	response.Response
}

//go:generate mockery
type UserSaver interface {
	SaveUser(ctx context.Context, username string) (int64, error)
}

func New(userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)
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

			errMsg := "invalid request"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.ValidationError(validateErr)})
			return
		}

		userId, err := userSaver.SaveUser(r.Context(), req.Username)
		if err != nil {
			errMsg := "save user failed"
			slog.Error(errMsg, sl.Err(err))
			render.JSON(w, r, Response{Response: response.Error(errMsg)})
			return
		}

		render.JSON(w, r, Response{
			Id:       userId,
			Username: req.Username,
			Response: response.Ok(),
		})
	}
}
