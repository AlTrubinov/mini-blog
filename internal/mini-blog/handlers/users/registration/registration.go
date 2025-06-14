package registration

import (
	"context"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/lib/auth"
	"mini-blog/pkg/apperror"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Id int64 `json:"id,omitempty"`
	auth.JwtToken
	response.Response
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, passwordHash string) (int64, error)
}

func New(userSaver UserSaver, tm *auth.TokenManager) http.HandlerFunc {
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

		slog.Info("user registration request body decoded", slog.Any("request_username", req.Username))

		if err := validapi.Request(req); err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
		}

		userId, err := userSaver.SaveUser(r.Context(), req.Username, string(hash))
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		token, err := tm.Generate(userId)
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
			JwtToken: token,
			Response: response.Created(),
		})
	}
}
