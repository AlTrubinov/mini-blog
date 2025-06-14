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

func NotFoundError(msg string) Response {
	return Response{
		Code:  http.StatusNotFound,
		Error: msg,
	}
}

func ForbiddenError(msg string) Response {
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

func TimeoutError(msg string) Response {
	return Response{
		Code:  http.StatusRequestTimeout,
		Error: msg,
	}
}

func UnauthorizedError(msg string) Response {
	return Response{
		Code:  http.StatusUnauthorized,
		Error: msg,
	}
}

func Json(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	w.WriteHeader(code)
	render.JSON(w, r, v)
}

func GetErrorResponseByCode(code int, msg string) Response {
	switch code {
	case http.StatusBadRequest:
		return ValidationError(msg)
	case http.StatusNotFound:
		return NotFoundError(msg)
	case http.StatusForbidden:
		return ForbiddenError(msg)
	case http.StatusRequestTimeout:
		return TimeoutError(msg)
	case http.StatusUnauthorized:
		return UnauthorizedError(msg)
	default:
		return InternalError(msg)
	}
}
