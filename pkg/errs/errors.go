package errs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNotFound     = errors.New("Not Found")
	ErrBadRequest   = errors.New("Bad Request")
	ErrExists       = errors.New("Exists")
	ErrConflict     = errors.New("Conflict")
	ErrUnauthorized = errors.New("Unauthorized")
)

func BadRequest(s string, p ...interface{}) error {
	return fmt.Errorf("[400 Bad Request] %s", fmt.Sprintf(s, p...))
}

func NotFound(s string, p ...interface{}) error {
	return fmt.Errorf("[404 Not Found] %s", fmt.Sprintf(s, p...))
}

func Unauthorized(s string, p ...interface{}) error {
	return fmt.Errorf("[401 Unauthorized] %s", fmt.Sprintf(s, p...))
}

//type welawError struct {
//Code    int    `json:"code"`
//Message string `json:"message"`
//}

//func (err *welawError) Error() string {
//return fmt.Sprintf("[%d] %s", err.Code, err.Message)
//}

//func NewError(message string) error {
//return &welawError{Code: 0, Message: message}
//}

func ErrToStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func StrToErr(s string) error {
	switch {
	case s == "":
		return nil
	case s == ErrNotFound.Error(), strings.HasPrefix(s, "[404 Not Found]"):
		return ErrNotFound
	case s == ErrBadRequest.Error(), strings.HasPrefix(s, "[400 Bad Request]"):
		return ErrBadRequest
	case s == ErrUnauthorized.Error(), strings.HasPrefix(s, "[401 Unauthorized]"):
		return ErrUnauthorized
	case s == ErrConflict.Error():
		return ErrConflict
	}
	return errors.New(s)
}

func errToHttpCode(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrConflict:
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(errToHttpCode(err))
	json.NewEncoder(w).Encode(ErrorWrapper{Error: err.Error()})
}

func ErrorDecoder(r *http.Response) error {
	var w ErrorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type ErrorWrapper struct {
	Error string `json:"error"`
}
