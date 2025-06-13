package validapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Int64UrlParam(r *http.Request, paramName string) (int64, error) {
	urlParamStr := chi.URLParam(r, paramName)
	urlParam, err := strconv.ParseInt(urlParamStr, 10, 64)
	if err != nil {
		errMsg := "invalid %s: %s"
		return 0, fmt.Errorf(errMsg, paramName, urlParamStr)
	}
	return urlParam, nil
}

func JsonBodyDecode(r *http.Request, v interface{}) error {
	err := render.DecodeJSON(r.Body, &v)
	if err != nil {
		errMsg := "invalid JSON"
		return fmt.Errorf(errMsg)
	}
	return nil
}

func Request(v interface{}) error {
	if err := validate.Struct(v); err != nil {
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			return fmt.Errorf(validateErrMsg(validateErr))
		}
		errMsg := "request validation failed by internal error: %w"
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

func validateErrMsg(errs validator.ValidationErrors) string {
	var errMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is required", err.Field()))
		case "url":
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is invalid url", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is not valid", err.Field()))
		}
	}

	return strings.Join(errMessages, ", ")
}
