package registration

import (
	"context"
	"log/slog"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/pkg/apperror"
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
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		slog.Info("request body decoded", slog.Any("request", req))

		if err := validapi.Request(req); err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		userId, err := userSaver.SaveUser(r.Context(), req.Username)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		response.Json(w, r, http.StatusCreated, Response{
			Id:       userId,
			Username: req.Username,
			Response: response.Created(),
		})
	}
}
