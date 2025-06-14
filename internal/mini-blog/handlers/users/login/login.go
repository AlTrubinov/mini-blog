package login

import (
	"context"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/lib/auth"
	"mini-blog/internal/models/user"
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

type UserGetter interface {
	GetUser(ctx context.Context, username string) (user.User, error)
}

func New(userGetter UserGetter, tm *auth.TokenManager) http.HandlerFunc {
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

		slog.Info("get user request body decoded", slog.Any("request_username", req.Username))

		if err := validapi.Request(req); err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		foundedUser, err := userGetter.GetUser(r.Context(), req.Username)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundedUser.Password), []byte(req.Password))
		if err != nil {
			errMsg := apperror.ErrUnauthorized.Error()
			slog.Error(errMsg)
			response.Json(w, r, 401, response.UnauthorizedError(errMsg))
			return
		}

		token, err := tm.Generate(foundedUser.Id)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		response.Json(w, r, http.StatusOK, Response{
			Id:       foundedUser.Id,
			JwtToken: token,
			Response: response.Ok(),
		})
	}
}
