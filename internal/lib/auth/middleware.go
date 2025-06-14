package auth

import (
	"context"
	"net/http"
	"strings"

	"mini-blog/internal/lib/api/response"
	"mini-blog/pkg/apperror"
)

func (tm *TokenManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Json(w, r, 401, response.UnauthorizedError("authorization header is missing"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Json(w, r, 401, response.UnauthorizedError("invalid authorization header"))
			return
		}

		token := parts[1]
		claims, err := tm.Parse(token)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			response.Json(w, r, 403, response.ForbiddenError("invalid token payload"))
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserId, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
