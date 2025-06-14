package create

import (
	"context"
	"log/slog"
	"mini-blog/internal/lib/auth"
	"net/http"

	"mini-blog/internal/lib/api/response"
	"mini-blog/internal/lib/api/validapi"
	"mini-blog/pkg/apperror"
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

func New(creator NotesCreator, tm *auth.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		userId, err := validapi.Int64UrlParam(r, "user_id")
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		if err = tm.CheckUserAccess(r.Context(), userId); err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		err = validapi.JsonBodyDecode(r, &req)
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

		noteId, err := creator.CreateNote(r.Context(), userId, req.Title, req.Content)
		if err != nil {
			errMsg := err.Error()
			errCode := apperror.GetCodeByError(err)
			errResp := response.GetErrorResponseByCode(errCode, errMsg)
			slog.Error(errMsg)
			response.Json(w, r, errCode, errResp)
			return
		}

		response.Json(w, r, http.StatusCreated, Response{
			Id:       noteId,
			Response: response.Created(),
		})
	}
}
