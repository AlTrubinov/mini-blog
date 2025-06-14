package token

import (
	"log/slog"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/internal/lib/auth"
	"mini-blog/pkg/apperror"
)

type Request struct {
	RefreshToken string `json:"refresh" validate:"required"`
}

type Response struct {
	auth.JwtToken
	response.Response
}

func New(tm *auth.TokenManager) http.HandlerFunc {
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

		slog.Info("refresh token request body decoded")

		if err := validapi.Request(req); err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		claims, err := tm.Parse(req.RefreshToken)
		if err != nil {
			errMsg := apperror.ErrUnauthorized.Error()
			slog.Error(errMsg)
			response.Json(w, r, 401, response.UnauthorizedError(errMsg))
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			response.Json(w, r, 403, response.ForbiddenError("invalid token payload"))
			return
		}

		token, err := tm.Generate(int64(userID))
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		response.Json(w, r, http.StatusOK, Response{
			JwtToken: token,
			Response: response.Ok(),
		})
	}
}
