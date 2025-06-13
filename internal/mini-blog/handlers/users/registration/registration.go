package registration

import (
	"context"
	"log/slog"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
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

		err := validapi.JsonBodyDecode(r, &req)
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

		userId, err := userSaver.SaveUser(r.Context(), req.Username)
		if err != nil {
			errMsg := "save user failed"
			slog.Error(errMsg, sl.Err(err))
			response.Json(w, r, http.StatusBadRequest, response.ValidationError(err.Error()))
			return
		}

		response.Json(w, r, http.StatusCreated, Response{
			Id:       userId,
			Username: req.Username,
			Response: response.Created(),
		})
	}
}
