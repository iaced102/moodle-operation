package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	switch e.Code {
	case "not_found":
		return http.StatusNotFound
	case "invalid_input":
		return http.StatusBadRequest
	case "conflict":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func  Conflict(message string, err error) *AppError {
	return New("conflict", message, err)
}

func Internal(message string, err error) *AppError {
	return New("internal", message, err)
}

func InvalidInput(message string, err error) *AppError {
	return New("invalid_input", message, err)
}
