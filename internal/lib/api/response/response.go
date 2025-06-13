package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}

func Ok() Response {
	return Response{
		Code: http.StatusOK,
	}
}

func Created() Response {
	return Response{
		Code: http.StatusCreated,
	}
}

func ValidationError(msg string) Response {
	return Response{
		Code:  http.StatusBadRequest,
		Error: msg,
	}
}

func NotFound(msg string) Response {
	return Response{
		Code:  http.StatusNotFound,
		Error: msg,
	}
}

func Forbidden(msg string) Response {
	return Response{
		Code:  http.StatusForbidden,
		Error: msg,
	}
}

func InternalError(msg string) Response {
	return Response{
		Code:  http.StatusInternalServerError,
		Error: msg,
	}
}

func Json(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	w.WriteHeader(code)
	render.JSON(w, r, v)
}
